package gateway

import (
	"database/sql"

	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/rdcloud-io/global"
	"github.com/rdcloud-io/global/servicelist"
	"github.com/rdcloud-io/openapi/common"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var Conf *common.Config

func Start() {
	common.InitConfig()
	Conf = common.Conf
	global.InitLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug, servicelist.OpenapiGateway)
	Logger = global.Logger

	db = global.InitMysql(&global.MysqlConfig{
		Acc:      Conf.Mysql.Acc,
		Pw:       Conf.Mysql.Pw,
		Addr:     Conf.Mysql.Addr,
		Port:     Conf.Mysql.Port,
		Database: Conf.Mysql.Database,
	})
	initEtcd()
	initUpdateApi()

	go watchUpstramServers()

	apis = &Apis{
		&sync.Map{},
	}
	apis.LoadAll()

	e := echo.New()
	e.Any("/*", apiRoute)
	e.Logger.Fatal(e.Start(":" + Conf.Api.GatewayPort))
}

var db *sql.DB
