package gateway

import (
	"net/http"
	"time"

	"fmt"

	"github.com/labstack/echo"
	"github.com/rdcloud-io/openapi/global"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var cli = &fasthttp.Client{}

func apiRoute(c echo.Context) error {
	ts := time.Now()

	// 生成request_id
	rid := global.RequestID()
	// 获取debug选项
	if c.FormValue("log_debug") == "on" {
		c.Set("debug_on", true)
	} else {
		c.Set("debug_on", false)
	}

	apiName := c.FormValue("api_name")

	// 查询是否存在此Api
	apiI, ok := apis.Load(apiName)
	// 记录请求IP、参数
	Logger.Info("收到新请求", zap.String("rid", rid), zap.String("ip", c.RealIP()), zap.String("api_name", apiName), zap.Bool("api_exist", ok))
	if !ok {
		// api不存在，返回错误
		return c.String(http.StatusOK, "Api不存在")
	}

	api := apiI.(*Api)

	global.DebugLog(rid, c.Get("debug_on").(bool), "请求的api", zap.Any("api", *api))

	// 生成url
	var upstreamUrl string

	if len(api.UpstreamServers) <= 0 {
		Logger.Warn("要请求的服务不存活", zap.String("rid", rid), zap.String("api_name", apiName))
		return c.String(http.StatusOK, "该Api没有服务器存活")
	}
	if api.ProxyMode == "1" {
		uri := c.Request().RequestURI
		//  api.UpstreamMode 为 1和2时，处理方式暂时统一
		url := api.UpstreamServers[0].IP
		upstreamUrl = url + uri
	} else {
		upstreamUrl = api.UpstreamServers[0].IP
	}

	global.DebugLog(rid, c.Get("debug_on").(bool), "raw upstream url", zap.String("uri", c.Request().RequestURI), zap.String("raw_upstream_url", upstreamUrl))
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	// 设置Method
	req.Header.SetMethod(api.Method)

	// 透传参数
	args := &fasthttp.Args{}
	params, err := c.FormParams()
	if err != nil {
		Logger.Info("解析请求参数错误", zap.String("rid", rid), zap.Error(err), zap.String("paramas", params.Encode()), zap.String("api_name", apiName))
		return c.String(http.StatusOK, "获取参数错误")
	}
	for param := range params {
		args.Set(param, c.FormValue(param))
	}
	args.Set("rid", rid)
	global.DebugLog(rid, c.Get("debug_on").(bool), "请求参数", zap.String("paramas", params.Encode()))
	// args.Del("api_name")

	// 透传cookie
	for _, cookie := range c.Cookies() {
		// global.DebugLog(c.Get("rid").(string), c.Get("debug_on").(bool), zap.String("cookie_name", cookie.Name), zap.String("cookie_val", cookie.Value))
		req.Header.SetCookie(cookie.Name, cookie.Value)
	}

	global.DebugLog(rid, c.Get("debug_on").(bool), "请求cookies", zap.Any("cookie", c.Cookies()))
	// 设置X-FORWARD-FOR
	req.Header.Set("X-Forwarded-For", c.RealIP())

	// 请求
	switch api.Method {
	case "POST":
		args.WriteTo(req.BodyWriter())
	case "GET":
		// 拼接url
		upstreamUrl = upstreamUrl + "?" + args.String()

	}

	global.DebugLog(rid, c.Get("debug_on").(bool), "最终upstream url", zap.String("upstream_url", upstreamUrl))
	req.SetRequestURI(upstreamUrl)

	err = cli.DoTimeout(req, resp, 10*time.Second)

	if err != nil {
		Logger.Info("api请求错误", zap.String("rid", rid), zap.Error(err), zap.String("api_name", apiName))
		return c.String(resp.StatusCode(), err.Error())
	}

	if resp.StatusCode() != 200 {
		Logger.Info("api请求code不为200", zap.String("rid", rid), zap.Int("code", resp.StatusCode()), zap.String("api_name", apiName))
		return c.String(resp.StatusCode(), fmt.Sprintf("请求返回Code异常：%v", resp.StatusCode()))
	}

	Logger.Info("api请求成功", zap.String("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000), zap.String("api_name", apiName))
	global.DebugLog(rid, c.Get("debug_on").(bool), "api请求返回body", zap.String("body", string(resp.Body())))

	return c.String(http.StatusOK, string(resp.Body()))
}
