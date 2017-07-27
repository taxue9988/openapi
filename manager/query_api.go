package manager

import (
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/rdcloud-io/global/apilist"
	"github.com/rdcloud-io/openapi/common"
	"github.com/rdcloud-io/sdk/api"
	"github.com/rdcloud-io/sdk/etcd"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var etcdCli *clientv3.Client

var Servers = &sync.Map{}

var cli = &fasthttp.Client{}

func initGatewayUpdate() {
	go func() {
		etcdCli = etcd.Init(common.Conf.Etcd.Addrs, Logger)

		resCh, errCh := api.QueryServerByAPI(etcdCli, []string{apilist.OpenapiGatewayUpdateApi, apilist.StaffCheckLogin}, 0)
		for {
			select {
			case servers := <-resCh:
				Servers.Store(servers.ApiName, servers)
				Logger.Debug("更新服务器列表", zap.Any("gateway_addr", servers))
			case err := <-errCh:
				Logger.Warn("请求etcd异常", zap.Error(err), zap.Any("etcd_addr", common.Conf.Etcd.Addrs))
			}
		}
	}()

}

func updateApi(apiName string, tp int) {
	args := &fasthttp.Args{}
	args.Set("api_name", apiName)
	args.SetUint("type", tp)

	serversS, ok := Servers.Load(apilist.OpenapiGatewayUpdateApi)
	if !ok {
		Logger.Warn("gateway没有节点存活")
		return
	}
	servers := serversS.(*api.QueryServerRes)

	for _, server := range servers.Servers {
		url := "http://" + server.IP + server.Path
		code, _, err := cli.Post(nil, url, args)
		if err != nil {
			Logger.Warn("manager update api error", zap.Error(err), zap.String("url", url), zap.String("args", args.String()))
			continue
		}

		if code != 200 {
			Logger.Warn("manager update api code invalid", zap.Int("code", code), zap.String("url", url), zap.String("args", args.String()))
			continue
		}
	}
}
