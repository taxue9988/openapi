package manager

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"fmt"

	"github.com/labstack/echo"
	"github.com/openapm/talents"
	"github.com/rdcloud-io/global"
	"github.com/rdcloud-io/global/code"
	"github.com/rdcloud-io/openapi/data"
)

func inApiCreate(c echo.Context) error {
	ts := time.Now()
	api := &data.InApi{}

	api.Company = global.FormValueTrimed("company", c)
	api.Product = global.FormValueTrimed("product", c)
	api.System = global.FormValueTrimed("system", c)
	api.Interface = global.FormValueTrimed("interface", c)

	if api.Company == "" || api.Product == "" || api.System == "" || api.Interface == "" {
		return c.JSON(http.StatusOK, &global.Result{
			Code:    code.InApiParamCantBeEmpty,
			Message: "必填项不能为空",
		})
	}

	api.ApiName = api.Company + "." + api.Product + "." + api.System + "." + api.Interface

	api.Method = c.FormValue("method")
	api.Register = global.FormValueTrimed("register", c)
	api.ApiDesc = c.FormValue("api_desc")
	api.ParamDesc = c.FormValue("param_desc")
	api.ReturnDesc = global.FormValueTrimed("return_desc", c)

	api.InputDate = talents.Time2String(time.Now())
	staff := c.Get("staff").(map[string]interface{})
	api.InputStaff = staff["name"].(string) + "/" + staff["uid"].(string)

	query := fmt.Sprintf("INSERT INTO in_api (`api_name`,`company`,`product`,`system`,`interface`,`method`,`register`,`api_desc`,`param_desc`,`return_desc`,`input_date`,`input_staff`)VALUES ('%s', '%s','%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')", api.ApiName, api.Company, api.Product, api.System, api.Interface, api.Method, api.Register, api.ApiDesc,
		api.ParamDesc, api.ReturnDesc, api.InputDate, api.InputStaff)
	_, err := db.Exec(query)
	if err != nil {
		Logger.Info("api create, insert error", zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "create api error")
	}

	Logger.Info("api创建成功", zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))

	return c.JSON(http.StatusOK, &global.Result{
		Code: code.Success,
		Data: map[string]interface{}{
			"api": api,
		},
	})
}

func inApiList(c echo.Context) error {
	ts := time.Now()

	query := fmt.Sprintf("SELECT * FROM in_api")
	rows, err := db.Query(query)
	if err != nil {
		Logger.Info("api query error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, global.Result{
			Code:    code.InApiListMysqlQuerError,
			Message: "获取列表错误",
		})
	}
	defer rows.Close()

	var apis []*data.InApi
	for rows.Next() {
		tempApi := &data.InApi{}
		err := rows.Scan(&tempApi.ApiName, &tempApi.Company, &tempApi.Product,
			&tempApi.System, &tempApi.Interface, &tempApi.Method, &tempApi.Register, &tempApi.ApiDesc, &tempApi.ParamDesc,
			&tempApi.ReturnDesc, &tempApi.InputDate, &tempApi.InputStaff)
		if err != nil {
			Logger.Info("api scan error", zap.Error(err), zap.String("query", query))
			return c.JSON(http.StatusOK, global.Result{
				Code:    code.InApiListMysqlScanError,
				Message: "获取列表错误",
			})
		}
		apis = append(apis, tempApi)
	}

	Logger.Info("api列表查询成功", zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))
	if len(apis) == 0 {
		return c.JSON(http.StatusOK, &global.Result{
			Code:    code.InApiListQueryEmpty,
			Message: "还没创建Api",
		})
	}
	return c.JSON(http.StatusOK, &global.Result{
		Code: code.Success,
		Data: map[string]interface{}{
			"apis": apis,
		},
	})
}

func inApiDelete(c echo.Context) error {
	ts := time.Now()

	apiAll := c.FormValue("apis")
	Logger.Info("api删除", zap.String("api_name", apiAll))

	apis := strings.Split(apiAll, ",")
	for _, api := range apis {
		query := fmt.Sprintf("DELETE FROM in_api where `api_name`='%s'", api)
		_, err := db.Exec(query)
		if err != nil {
			Logger.Info("api delete  error", zap.Error(err), zap.String("query", query))
			return c.JSON(http.StatusOK, global.Result{
				Code:    code.InApiDeleteMysqlQuerError,
				Message: "删除api错误",
			})
		}
	}

	Logger.Info("api删除成功", zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))
	return c.JSON(http.StatusOK, global.Result{
		Code:    code.Success,
		Message: "删除成功",
	})
}

func inApiUpdate(c echo.Context) error {
	ts := time.Now()

	api := &data.InApi{}
	api.ApiName = c.FormValue("api_name")

	api.Method = c.FormValue("method")
	api.Register = global.FormValueTrimed("register", c)
	api.ReturnDesc = global.FormValueTrimed("return_desc", c)
	api.ApiDesc = c.FormValue("api_desc")
	api.ParamDesc = c.FormValue("param_desc")

	if api.ApiName == "" {
		Logger.Info("api_name不能为空")
		return c.JSON(http.StatusOK, &global.Result{
			Code:    code.InApiUpdateApiNameEmpty,
			Message: "Api名不能为空",
		})
	}

	Logger.Info("api更新", zap.Any("api", *api))

	updateDate := talents.Time2String(time.Now())
	staff := c.Get("staff").(map[string]interface{})
	updateStaff := staff["name"].(string) + "/" + staff["uid"].(string)

	query := fmt.Sprintf("UPDATE in_api SET `method`='%s',`api_desc`='%s',`param_desc`='%s' WHERE `api_name`='%s'",
		api.Method, api.ApiDesc, api.ParamDesc, api.ApiName)
	_, err := db.Exec(query)
	if err != nil {
		Logger.Info("api update  error", zap.Error(err), zap.String("query", query))
		return c.JSON(http.StatusOK, &global.Result{
			Code:    code.InApiUpdateMysqlQuerError,
			Message: "Api更新错误",
		})
	}

	// 更新revision
	query = fmt.Sprintf("INSERT INTO in_api_revision (`api_name`,`method`,`register`,`api_desc`,`param_desc`,`return_desc`,`update_date`,`update_staff`)VALUES ('%s','%s','%s','%s','%s','%s','%s','%s')", api.ApiName, api.Method, api.Register, api.ApiDesc,
		api.ParamDesc, api.ReturnDesc, updateDate, updateStaff)
	_, err = db.Exec(query)
	if err != nil {
		Logger.Info("api revision error", zap.Error(err), zap.String("query", query))
	}

	Logger.Info("api更新成功", zap.Int64("eclapsed", time.Now().Sub(ts).Nanoseconds()/1000))

	return c.JSON(http.StatusOK, &global.Result{
		Code:    code.Success,
		Message: "更新Api成功",
	})
}
