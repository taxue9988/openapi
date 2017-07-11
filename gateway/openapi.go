package gateway

import (
	"database/sql"
	"fmt"

	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/rdcloud-io/openapi/global"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var Conf *global.Config

func Start() {
	global.InitConfig()
	Conf = global.Conf
	global.InitLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
	Logger = global.Logger

	initMysql()
	initEtcd()
	go watchUpstramServers()

	apis = &Apis{
		&sync.Map{},
	}
	apis.LoadAll()

	e := echo.New()
	e.Any("/*", apiRoute)
	e.Logger.Fatal(e.Start(Conf.Api.Addr))
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
