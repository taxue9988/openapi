package service

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

/* 从Mysql中加载API信息到内存中*/
type ApiRaw struct {
	ID            int
	Domain        string
	Group         string
	Name          string
	Version       string
	Method        string
	ProxyMode     int
	UpstreamMode  int
	UpstreamValue string
}

func loadApis() {
	query := fmt.Sprintf("SELECT * FROM api")
	rows, err := db.Query(query)
	if err != nil {
		Logger.Fatal("query openapi.api error ", zap.Error(err))
	}

	for rows.Next() {
		rawApi := &ApiRaw{}
		err = rows.Scan(&rawApi.ID, &rawApi.Domain, &rawApi.Group, &rawApi.Name, &rawApi.Version,
			&rawApi.Method, &rawApi.ProxyMode, &rawApi.UpstreamMode, &rawApi.UpstreamValue)
		if err != nil {
			Logger.Fatal("scan openapi.api error ", zap.Error(err))
		}

		api := &Api{}
		api.Method = rawApi.Method
		api.ProxyMode = rawApi.ProxyMode
		api.UpstreamMode = rawApi.UpstreamMode
		api.Name = rawApi.Domain + "." + rawApi.Group + "." + rawApi.Name + "." + rawApi.Version
		if api.UpstreamMode == 1 {
			api.UpstreamServers = []*UpstreamServer{
				&UpstreamServer{
					IP:   rawApi.UpstreamValue,
					Load: 1,
				},
			}
		} else {
			// 从etcd读取key=api.Name的值
		}

		Apis.Store(api.Name, api)
	}

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

func mockApis() {
	// 支付收银台结算
	payCheckApi := &Api{
		Name:         "pay.cashier.check.v1",
		Method:       "POST",
		ProxyMode:    1,
		UpstreamMode: 1,
		UpstreamServers: []*UpstreamServer{
			&UpstreamServer{
				IP:   "http://localhost:1324",
				Load: 1,
			},
		},
	}

	Apis.Store(payCheckApi.Name, payCheckApi)

	payCheckApiV2 := &Api{
		Name:         "pay.cashier.check.v2",
		Method:       "POST",
		ProxyMode:    1,
		UpstreamMode: 1,
		UpstreamServers: []*UpstreamServer{
			&UpstreamServer{
				IP:   "http://localhost:1324/a",
				Load: 1,
			},
		},
	}

	Apis.Store(payCheckApiV2.Name, payCheckApiV2)

	// 支付收银台查询
	payQueryApi := &Api{
		Name:         "pay.cashier.query.v1",
		Method:       "GET",
		ProxyMode:    2,
		UpstreamMode: 1,
		UpstreamServers: []*UpstreamServer{
			&UpstreamServer{
				IP:   "http://localhost:1325/pay/cashier/query",
				Load: 1,
			},
		},
	}

	Apis.Store(payQueryApi.Name, payQueryApi)

	// HTTP信息查询
	httpGetApi := &Api{
		Name:         "http.info.get.v1",
		Method:       "GET",
		ProxyMode:    1,
		UpstreamMode: 1,
		UpstreamServers: []*UpstreamServer{
			&UpstreamServer{
				IP:   "http://httpbin.org",
				Load: 1,
			},
		},
	}

	Apis.Store(httpGetApi.Name, httpGetApi)

	httpGetApiV2 := &Api{
		Name:         "http.info.get.v2",
		Method:       "GET",
		ProxyMode:    2,
		UpstreamMode: 1,
		UpstreamServers: []*UpstreamServer{
			&UpstreamServer{
				IP:   "http://httpbin.org/get",
				Load: 1,
			},
		},
	}

	Apis.Store(httpGetApiV2.Name, httpGetApiV2)
}
