package manager

import (
	"database/sql"

	"github.com/labstack/echo"
	"github.com/rdcloud-io/global"
	"github.com/rdcloud-io/global/code"
	"github.com/rdcloud-io/openapi/common"

	"github.com/rdcloud-io/global/apilist"
	"github.com/rdcloud-io/sdk/account"
	"github.com/rdcloud-io/sdk/api"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var Conf *common.Config

var db *sql.DB

func Start() {
	common.InitConfig()
	Conf = common.Conf

	global.InitLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug, Conf.Common.Service)
	Logger = global.Logger

	db = global.InitMysql(&global.MysqlConfig{
		Acc:      Conf.Mysql.Acc,
		Pw:       Conf.Mysql.Pw,
		Addr:     Conf.Mysql.Addr,
		Port:     Conf.Mysql.Port,
		Database: Conf.Mysql.Database,
	})
	initGatewayUpdate()

	e := echo.New()
	e.Static("/", "manager/public/docs")
	e.POST("/api/create", apiCreate)
	e.POST("/api/update", apiUpdate)
	e.POST("/api/query", apiQuery)
	e.POST("/api/delete", apiDelete)
	e.GET("/api/list", apiList, storeStaff, account.CheckStaffLogin)

	// 内部API管理
	e.POST("/inapi/create", inApiCreate, storeStaff, account.CheckStaffLogin)
	e.GET("/inapi/list", inApiList, storeStaff, account.CheckStaffLogin)
	e.POST("/inapi/delete", inApiDelete, storeStaff, account.CheckStaffLogin)
	e.POST("/inapi/update", inApiUpdate, storeStaff, account.CheckStaffLogin)
	e.Logger.Fatal(e.Start(":" + Conf.Admin.ManagerPort))
}

// 存储staff ip、port、路径
func storeStaff(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		global.SetCrossDomain(c)
		serversS, ok := Servers.Load(apilist.StaffCheckLogin)
		if !ok {
			return c.JSON(200, global.Result{
				Code:      code.StaffServersEmpty,
				Message:   "请重新登录",
				NeedLogin: true,
			})
		}

		servers := serversS.(*api.QueryServerRes)
		c.Set("addr", "http://"+servers.Servers[0].IP+servers.Servers[0].Path)

		return f(c)
	}
}
