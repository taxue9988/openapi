package gateway

/* 定期将自身状况报告给manager */

func initReport() {
	// 启动一个goroutine更新manager的地址
	go func() {
		// res, err := etcdCli.Get(context.Background(), global.Conf.Etcd.ServerKey+"rdcloud.openapi.manager", clientv3.WithPrefix())
	}()
	// 启动一个goroutine定期上报自身地址
}
