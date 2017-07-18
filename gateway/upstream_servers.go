package gateway

import (
	"bytes"
	"context"
	"sort"
	"strconv"

	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/rdcloud-io/global/apilist"
	"go.uber.org/zap"
)

func watchUpstramServers() {
	for {
		rch := etcdCli.Watch(context.Background(), Conf.Etcd.ServerKey, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == 0 { // put
					// 解析load
					load, _ := strconv.ParseFloat(string(ev.Kv.Value), 64)

					// 解析出ip
					ipIndex := bytes.LastIndex(ev.Kv.Key, []byte{'/'})
					ip := "http://" + string(ev.Kv.Key[ipIndex+1:])

					rest := ev.Kv.Key[:ipIndex]

					// 解析出apiName
					apiIndex := bytes.LastIndex(rest, []byte{'/'})
					apiName := string(rest[apiIndex+1:])

					// 如果api是gateway自身的api,则略过
					if apiName == apilist.OpenapiGatewayUpdateApi {
						continue
					}

					apiI, ok := apis.Load(apiName)
					if !ok {
						// api不存在 返回错误
						Logger.Info("api不存在，但是取到了服务器的打点信息", zap.String("api", apiName))
						continue
					}

					api := apiI.(*Api)

					// 更新对应的ip
					ipExist := false
					for _, server := range api.UpstreamServers {
						if server.IP == ip {
							server.Load = load
							ipExist = true
						}
					}

					// 不存在该ip,则添加新的ip信息
					if !ipExist {
						newServer := &UpstreamServer{
							IP:   ip,
							Load: load,
						}
						api.UpstreamServers = append(api.UpstreamServers, newServer)
					}

					sort.Slice(api.UpstreamServers, func(i, j int) bool {
						return api.UpstreamServers[i].Load < api.UpstreamServers[j].Load
					})

					fmt.Printf("watch, key: %s, ip : %s, load: %v\n", apiName, ip, load)
					for _, s := range api.UpstreamServers {
						fmt.Println("etcd watch插入,该api的最新服务器列表: ", *s)
					}
				} else if ev.Type == 1 { // delete
					// 解析出ip
					ipIndex := bytes.LastIndex(ev.Kv.Key, []byte{'/'})
					ip := "http://" + string(ev.Kv.Key[ipIndex+1:])

					rest := ev.Kv.Key[:ipIndex]

					// 解析出apiName
					apiIndex := bytes.LastIndex(rest, []byte{'/'})
					apiName := string(rest[apiIndex+1:])

					apiI, ok := apis.Load(apiName)
					if !ok {
						// api不存在，返回错误
						Logger.Info("api不存在，但是取到了服务器的打点信息", zap.String("api", apiName))
						continue
					}

					api := apiI.(*Api)

					for i, server := range api.UpstreamServers {
						if i == len(api.UpstreamServers) {
							api.UpstreamServers = api.UpstreamServers[:i]
							break
						}

						if server.IP == ip {
							// 删除该server
							api.UpstreamServers = append(api.UpstreamServers[:i], api.UpstreamServers[i+1:]...)
						}
					}

					for _, s := range api.UpstreamServers {
						fmt.Println("etcd watch 删除,该api的最新服务器列表: ", *s)
					}
				}
			}
		}
	}
}
