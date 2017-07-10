package service

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"

	"github.com/coreos/etcd/clientv3"
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
			resp, err := etcdCli.Get(context.Background(), Conf.Etcd.ServerKey+api.Name, clientv3.WithPrefix())
			if err != nil {
				Logger.Fatal("etcd get error", zap.Error(err))
			}

			servers := make([]*UpstreamServer, 0, len(resp.Kvs))
			for _, v := range resp.Kvs {
				ip, load := ipAndLoad(v.Key, v.Value)
				servers = append(servers, &UpstreamServer{
					IP:   ip,
					Load: load,
				})
			}

			// 对负载进行从小到大的排列
			sort.Slice(servers, func(i, j int) bool {
				return servers[i].Load < servers[j].Load
			})

			api.UpstreamServers = servers

			for _, s := range api.UpstreamServers {
				fmt.Printf("api load: %s 的最新服务器列表: %v\n", api.Name, *s)
			}
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

func ipAndLoad(key []byte, val []byte) (string, float64) {
	// 解析load
	load, _ := strconv.ParseFloat(string(val), 64)

	// 解析出ip
	ipIndex := bytes.LastIndex(key, []byte{'/'})
	ip := "http://" + string(key[ipIndex+1:])

	return ip, load
}
