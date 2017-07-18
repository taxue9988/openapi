package manager

import (
	"database/sql"
	"fmt"

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

	initMysql()
	initGatewayUpdate()

	e := echo.New()
	e.Static("/", "manager/public/docs")
	e.POST("/api/create", apiCreate)
	e.POST("/api/update", apiUpdate)
	e.POST("/api/query", apiQuery)
	e.POST("/api/delete", apiDelete)
	e.GET("/api/list", apiList)
	e.Logger.Fatal(e.Start(":" + Conf.Admin.ManagerPort))
}

var db *sql.DB

func initMysql() {
	var err error

	// 初始化mysql连接
	sqlConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Conf.Mysql.Acc, Conf.Mysql.Pw,
		Conf.Mysql.Addr, Conf.Mysql.Port, Conf.Mysql.Database)
	db, err = sql.Open("mysql", sqlConn)
	if err != nil {
		Logger.Fatal("init mysql error", zap.Error(err))
	}

	// 测试db是否正常
	err = db.Ping()
	if err != nil {
		Logger.Fatal("init mysql, ping error", zap.Error(err))
	}
}
