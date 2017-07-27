package gateway

import (
	"strings"

	"fmt"

	"github.com/labstack/echo"
	"github.com/rdcloud-io/global"
	"github.com/rdcloud-io/global/apilist"
	"github.com/rdcloud-io/openapi/common"
	"github.com/rdcloud-io/sdk/api"
	"go.uber.org/zap"
)

func initUpdateApi() {
	go func() {
		e := echo.New()
		e.POST("/api/update", apiUpdate)
		e.Logger.Fatal(e.Start(":" + Conf.Api.ApiUpdatePort))
	}()

	go registerUpdateApi()
}
func registerUpdateApi() {
	servers := []*global.ServerInfo{
		&global.ServerInfo{
			APIName: apilist.OpenapiGatewayUpdateApi,
			IP:      Conf.Common.RealIp + ":" + Conf.Api.ApiUpdatePort,
			Path:    "/api/update",
			Load:    1,
		},
	}
	errCh := api.StoreServersByApi(etcdCli, servers)
	for {
		select {
		case err := <-errCh:
			Logger.Warn("请求etcd异常", zap.Error(err), zap.Any("etcd_addr", common.Conf.Etcd.Addrs))
		}
	}
}

func apiUpdate(c echo.Context) error {
	apiNames := strings.Split(c.FormValue("api_name"), ",")
	tp := c.FormValue("type")

	Logger.Info("api update", zap.Any("api_name", apiNames), zap.String("type", tp))
	for _, apiName := range apiNames {
		switch tp {
		case "1", "2": //创建api,更新api
			query := fmt.Sprintf("select * from api where full_name='%s'", apiName)
			rows, err := db.Query(query)
			if err != nil {
				Logger.Warn("query openapi.api error ", zap.Error(err))
			}
			for rows.Next() {
				load(rows)
			}
		case "3": //删除api
			apis.Delete(apiName)
		}
	}

	return nil
}
