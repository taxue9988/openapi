package manager

import (
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/labstack/echo"
	"github.com/rdcloud-io/openapi/apidata"
)

func apiCreate(c echo.Context) error {
	api, err := getApiPramas(c)
	if err != nil {
		Logger.Info("api parse error", zap.Error(err))
		return c.String(http.StatusOK, "api parse error")
	}
	api.FullName = api.Company + "." + api.Product + "." + api.System + "." + api.Interface + "." + api.Version
	query := fmt.Sprintf("INSERT INTO api (`full_name`,`company`,`product`,`system`,`interface`,`version`,`method`,`proxy_mode`,`upstream_mode`,`upstream_value`) VALUES ('%s', '%s','%s','%s','%s','%s','%s','%d','%d','%s')", api.FullName, api.Company, api.Product, api.System, api.Interface, api.Version, api.Method, api.ProxyMode,
		api.UpstreamMode, api.UpstreamValue)
	_, err = db.Exec(query)
	if err != nil {
		Logger.Info("api create, nsert error", zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "create api error")
	}

	return c.String(http.StatusOK, "success")
}

func apiUpdate(c echo.Context) error {
	api := &apidata.API{}
	api.FullName = c.FormValue("full_name")
	api.Method = c.FormValue("method")
	api.ProxyMode, _ = strconv.Atoi(c.FormValue("proxy_mode"))
	api.UpstreamMode, _ = strconv.Atoi(c.FormValue("upstream_mode"))
	api.UpstreamValue = c.FormValue("upstream_value")

	api, err := getApiPramas(c)
	if err != nil {
		Logger.Info("api parse error", zap.Error(err))
		return c.String(http.StatusOK, "api parse error")
	}

	query := fmt.Sprintf("UPDATE api SET `method`='%s',`proxy_mode`='%d',`upstream_mode`='%d',`upstream_value`='%s' WHERE `full_name`='%s'",
		api.Method, api.ProxyMode, api.UpstreamMode, api.UpstreamValue, api.FullName)
	_, err = db.Exec(query)
	if err != nil {
		Logger.Info("api update  error", zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "update api error")
	}

	return c.String(http.StatusOK, "success")
}

// 查询几种组合情况
// domain
// domain + group + name + version
func apiQuery(c echo.Context) error {
	tp := c.FormValue("query_type")

	var query string

	switch tp {
	case "full_name":
		apiAll := c.FormValue("full_name")

		query = fmt.Sprintf("SELECT * FROM api WHERE `full_name`='%s'", apiAll)
	}

	rows, err := db.Query(query)
	if err != nil {
		Logger.Info("api query error", zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "api query error")
	}
	defer rows.Close()

	var apis []*apidata.API
	for rows.Next() {
		tempApi := &apidata.API{}
		err := rows.Scan(&tempApi.ID, &tempApi.FullName, &tempApi.Company, &tempApi.Product,
			&tempApi.System, &tempApi.Interface, &tempApi.Version, &tempApi.Method, &tempApi.ProxyMode, &tempApi.UpstreamMode, &tempApi.UpstreamValue)
		if err != nil {
			Logger.Info("api scan error", zap.Error(err), zap.String("query", query))
			return c.String(http.StatusOK, "api scan error")
		}
		apis = append(apis, tempApi)
	}

	return c.JSON(http.StatusOK, apis)
}

func apiDelete(c echo.Context) error {
	apiAll := c.FormValue("full_name")

	query := fmt.Sprintf("DELETE FROM api where `full_name`='%s'", apiAll)
	_, err := db.Exec(query)
	if err != nil {
		Logger.Info("api delete  error", zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "delete api error")
	}

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
	api.ProxyMode, _ = strconv.Atoi(c.FormValue("proxy_mode"))
	api.UpstreamMode, _ = strconv.Atoi(c.FormValue("upstream_mode"))
	api.UpstreamValue = c.FormValue("upstream_value")

	return api, nil
}
