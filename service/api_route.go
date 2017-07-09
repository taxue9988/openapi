package service

import (
	"net/http"
	"time"

	"fmt"

	"github.com/labstack/echo"
	"github.com/valyala/fasthttp"
)

var cli = &fasthttp.Client{}

func apiRoute(c echo.Context) error {
	apiName := c.FormValue("api_name")

	// 查询是否存在此Api
	apiI, ok := Apis.Load(apiName)
	if !ok {
		// api不存在，返回错误
		return c.String(http.StatusOK, "Api不存在")
	}

	api := apiI.(*Api)

	// 生成url
	var upstreamUrl string
	if api.ProxyMode == 1 {
		uri := c.Request().RequestURI
		if api.UpstreamMode == 1 {
			url := api.UpstreamServers[0].IP
			upstreamUrl = url + uri
		}
	} else {
		upstreamUrl = api.UpstreamServers[0].IP
	}

	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	// 设置Method
	req.Header.SetMethod(api.Method)

	// 透传参数
	args := &fasthttp.Args{}
	params, err := c.FormParams()
	if err != nil {
		return c.String(http.StatusOK, "获取参数错误")
	}
	for param, _ := range params {
		args.Set(param, c.FormValue(param))
	}
	// args.Del("api_name")

	// 透传cookie
	for _, cookie := range c.Cookies() {
		req.Header.SetCookie(cookie.Name, cookie.Value)
	}

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

	println("url:", upstreamUrl)
	req.SetRequestURI(upstreamUrl)

	err = cli.DoTimeout(req, resp, 10*time.Second)

	if err != nil {
		return c.String(resp.StatusCode(), err.Error())
	}

	if resp.StatusCode() != 200 {
		return c.String(resp.StatusCode(), fmt.Sprintf("请求返回Code异常：%v", resp.StatusCode()))
	}
	return c.String(http.StatusOK, string(resp.Body()))

}
