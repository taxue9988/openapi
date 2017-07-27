package gateway

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/openapm/talents"
	"github.com/rdcloud-io/global"
	"github.com/rdcloud-io/openapi/data"
	"go.uber.org/zap"
)

type Api struct {
	// domain.group.service.version
	// 只能由字母、数字、点组成
	FullName string

	// GET POST
	Method string

	// 1.Raw: 将请求的Path直接append在upstream_value(url)后
	// 2.Indirect: 直接访问upstream_value(url)
	ProxyMode string

	// 1.直接寻址： url = upstream_value
	// 2.间接寻址: 在etcd中取出key为Api.Name的值，返回的数据结构存储在UpstreamValue
	UpstreamMode    string
	UpstreamServers []*UpstreamServer
}

type UpstreamServer struct {
	Load float64
	IP   string
}

type Apis struct {
	*sync.Map
}

var apis *Apis

func (a *Apis) LoadAll() {
	query := fmt.Sprintf("SELECT * FROM api")
	rows, err := db.Query(query)
	if err != nil {
		Logger.Fatal("query openapi.api error ", zap.Error(err))
	}
	defer rows.Close()

	for rows.Next() {
		load(rows)

	}

}

func load(rows *sql.Rows) {
	rawApi := &data.API{}
	err := rows.Scan(&rawApi.ID, &rawApi.FullName, &rawApi.Company, &rawApi.Product, &rawApi.System, &rawApi.Interface, &rawApi.Version,
		&rawApi.Method, &rawApi.ProxyMode, &rawApi.UpstreamMode, &rawApi.UpstreamValue)
	if err != nil {
		Logger.Fatal("scan openapi.api error ", zap.Error(err))
	}

	api := &Api{}
	api.Method = rawApi.Method
	api.ProxyMode = rawApi.ProxyMode
	api.UpstreamMode = rawApi.UpstreamMode
	api.FullName = rawApi.FullName
	if api.UpstreamMode == "1" {
		api.UpstreamServers = []*UpstreamServer{
			&UpstreamServer{
				IP:   rawApi.UpstreamValue,
				Load: 1,
			},
		}
	} else {
		// 从etcd读取key=api.Name的值
		resp, err := etcdCli.Get(context.Background(), global.APIsRootPath+api.FullName, clientv3.WithPrefix())
		if err != nil {
			Logger.Fatal("etcd get error", zap.Error(err))
		}

		servers := make([]*UpstreamServer, 0, len(resp.Kvs))
		for _, v := range resp.Kvs {
			ip, load, _ := ipAndLoad(v.Key, v.Value)
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
			fmt.Printf("api load: %s 的最新服务器列表: %v\n", api.FullName, *s)
		}
	}

	fmt.Println(*api)
	apis.Store(api.FullName, api)
}
func ipAndLoad(key []byte, val []byte) (string, float64, int) {
	// 解析load
	loadPath := strings.Split(talents.Bytes2String(val), "--")

	load, _ := strconv.ParseFloat(loadPath[0], 64)
	path := loadPath[1]

	// 解析出ip
	ipIndex := bytes.LastIndex(key, []byte{'/'})
	ip := "http://" + string(key[ipIndex+1:]) + path

	return ip, load, ipIndex
}
