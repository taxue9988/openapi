package manager

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"time"

	"strings"

	"github.com/labstack/echo"
	"github.com/rdcloud-io/openapi/apidata"
	"github.com/rdcloud-io/openapi/global"
)

type ExtApiRes struct {
	Suc  bool                   `json:"suc"`
	Data map[string]interface{} `json:"data"`
}

func apiCreate(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

	ts, rid, debugOn := global.LogExtra(c)

	api := &apidata.API{}
	// api, err := getApiPramas(c)
	// Logger.Info("api创建", zap.String("rid", rid), zap.Any("api", *api))
	// if err != nil {
	// 	Logger.Info("api parse error", zap.String("rid", rid), zap.Error(err))
	// 	return c.String(http.StatusOK, "api parse error")
	// }

	// if api.Company == "" || api.Product == "" || api.System == "" || api.Interface == "" || api.Version == "" {
	// 	return c.JSON(http.StatusOK, &ExtApiRes{
	// 		Suc: false,
	// 		Data: map[string]interface{}{
	// 			"msg": "必填项不能为空",
	// 		},
	// 	})
	// }
	// api.FullName = api.Company + "." + api.Product + "." + api.System + "." + api.Interface + "." + api.Version
	api.FullName = c.FormValue("api_name")
	if api.FullName == "" {
		return c.JSON(http.StatusOK, &ExtApiRes{
			Suc: false,
			Data: map[string]interface{}{
				"msg": "必填项不能为空",
			},
		})
	}
	api.Method = c.FormValue("method")
	api.UpstreamMode = c.FormValue("upstream_mode")
	api.UpstreamValue = c.FormValue("upstream_value")
	api.ProxyMode = "1"

	names := strings.Split(api.FullName, ".")
	if len(names) != 5 {
		return c.JSON(http.StatusOK, &ExtApiRes{
			Suc: false,
			Data: map[string]interface{}{
				"msg": "Api名必须为(公司.产品.系统.接口.版本号)格式，例如tf56.zeus.center.query.v1",
			},
		})
	}
	api.Company, api.Product, api.System, api.Interface, api.Version = names[0], names[1], names[2], names[3], names[4]
	query := fmt.Sprintf("INSERT INTO api (`full_name`,`company`,`product`,`system`,`interface`,`version`,`method`,`proxy_mode`,`upstream_mode`,`upstream_value`) VALUES ('%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')", api.FullName, api.Company, api.Product, api.System, api.Interface, api.Version, api.Method, api.ProxyMode,
		api.UpstreamMode, api.UpstreamValue)
	global.DebugLog(rid, debugOn, "创建api sql", zap.String("sql", query))
	_, err := db.Exec(query)
	if err != nil {
		Logger.Info("api create, insert error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "create api error")
	}

	Logger.Info("api创建成功", zap.String("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))

	return c.JSON(http.StatusOK, &ExtApiRes{
		Suc: true,
		Data: map[string]interface{}{
			"api": api,
		},
	})
}

func apiUpdate(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	ts, rid, debugOn := global.LogExtra(c)

	api := &apidata.API{}
	api.FullName = c.FormValue("api_name")
	api.Method = c.FormValue("method")
	api.ProxyMode = c.FormValue("proxy_mode")
	api.UpstreamMode = c.FormValue("upstream_mode")
	api.UpstreamValue = c.FormValue("upstream_value")

	if api.FullName == "" {
		Logger.Info("api_name不能为空", zap.String("rid", rid))
		return c.JSON(http.StatusOK, &ExtApiRes{
			Suc: false,
			Data: map[string]interface{}{
				"msg": "api名不能为空",
			},
		})
	}
	Logger.Info("api更新", zap.String("rid", rid), zap.Any("api", *api))

	query := fmt.Sprintf("UPDATE api SET `method`='%s',`proxy_mode`='%s',`upstream_mode`='%s',`upstream_value`='%s' WHERE `full_name`='%s'",
		api.Method, api.ProxyMode, api.UpstreamMode, api.UpstreamValue, api.FullName)
	global.DebugLog(rid, debugOn, "更新api sql", zap.String("sql", query))
	_, err := db.Exec(query)
	if err != nil {
		Logger.Info("api update  error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "update api error")
	}

	Logger.Info("api更新成功", zap.String("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))

	return c.JSON(http.StatusOK, &ExtApiRes{
		Suc: true,
	})
}

// 查询几种组合情况
// domain
// domain + group + name + version
func apiQuery(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	ts, rid, debugOn := global.LogExtra(c)

	var query string

	apiAll := c.FormValue("api_name")

	query = fmt.Sprintf("SELECT * FROM api WHERE `full_name`='%s'", apiAll)

	Logger.Info("api查询", zap.String("rid", rid), zap.String("query", query))

	rows, err := db.Query(query)
	if err != nil {
		Logger.Info("api query error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "api query error")
	}
	defer rows.Close()

	var api *apidata.API
	for rows.Next() {
		tempApi := &apidata.API{}
		err := rows.Scan(&tempApi.ID, &tempApi.FullName, &tempApi.Company, &tempApi.Product,
			&tempApi.System, &tempApi.Interface, &tempApi.Version, &tempApi.Method, &tempApi.ProxyMode, &tempApi.UpstreamMode, &tempApi.UpstreamValue)
		if err != nil {
			Logger.Info("api scan error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
			return c.String(http.StatusOK, "api scan error")
		}
		api = tempApi
	}

	if api == nil {
		return c.JSON(http.StatusOK, &ExtApiRes{
			Suc: false,
			Data: map[string]interface{}{
				"msg": "api不存在",
			},
		})
	}
	Logger.Info("api查询成功", zap.String("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))
	global.DebugLog(rid, debugOn, "查询api", zap.Any("api", api))

	return c.JSON(http.StatusOK, &ExtApiRes{
		Suc: true,
		Data: map[string]interface{}{
			"api": api,
		},
	})
}

func apiList(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

	ts, rid, debugOn := global.LogExtra(c)

	query := fmt.Sprintf("SELECT * FROM api")
	rows, err := db.Query(query)
	if err != nil {
		Logger.Info("api query error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "api query error")
	}
	defer rows.Close()

	var apis []*apidata.API
	for rows.Next() {
		tempApi := &apidata.API{}
		err := rows.Scan(&tempApi.ID, &tempApi.FullName, &tempApi.Company, &tempApi.Product,
			&tempApi.System, &tempApi.Interface, &tempApi.Version, &tempApi.Method, &tempApi.ProxyMode, &tempApi.UpstreamMode, &tempApi.UpstreamValue)
		if err != nil {
			Logger.Info("api scan error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
			return c.String(http.StatusOK, "api scan error")
		}
		apis = append(apis, tempApi)
	}

	Logger.Info("api列表查询成功", zap.String("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))
	global.DebugLog(rid, debugOn, "查询api", zap.Any("apis", apis))
	if len(apis) == 0 {
		return c.JSON(http.StatusOK, &ExtApiRes{
			Suc: false,
		})
	}
	return c.JSON(http.StatusOK, &ExtApiRes{
		Suc: true,
		Data: map[string]interface{}{
			"apis": apis,
		},
	})
}

func apiDelete(c echo.Context) error {
	c.Response().Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	ts, rid, debugOn := global.LogExtra(c)

	apiAll := c.FormValue("apis")
	Logger.Info("api查询", zap.String("rid", rid), zap.String("api_name", apiAll))

	apis := strings.Split(apiAll, ",")
	for _, api := range apis {
		query := fmt.Sprintf("DELETE FROM api where `full_name`='%s'", api)
		global.DebugLog(rid, debugOn, "删除api", zap.String("query", query))
		_, err := db.Exec(query)
		if err != nil {
			Logger.Info("api delete  error", zap.Error(err), zap.String("query", query), zap.String("rid", rid))
			return c.String(http.StatusOK, "delete api error")
		}
	}

	Logger.Info("api删除成功", zap.String("rid", rid), zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))

	return c.String(http.StatusOK, "success")
}

/*------------------------------ 内部API -----------------------------------------*/
func getApiPramas(c echo.Context) (*apidata.API, error) {
	api := &apidata.API{}
	api.Company = c.FormValue("company")
	api.Product = c.FormValue("product")
	api.System = c.FormValue("system")
	api.Interface = c.FormValue("interface")
	api.Version = c.FormValue("version")
	api.Method = c.FormValue("method")
	api.ProxyMode = c.FormValue("proxy_mode")
	api.UpstreamMode = c.FormValue("upstream_mode")
	api.UpstreamValue = c.FormValue("upstream_value")

	return api, nil
}
