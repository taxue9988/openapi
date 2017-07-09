package service

import (
	"sync"
)

type Api struct {
	// domain.group.service.version
	// 只能由字母、数字、点组成
	Name string

	// GET POST
	Method string

	// 1.Raw: 将请求的Path直接append在upstream_value(url)后
	// 2.Indirect: 直接访问upstream_value(url)
	ProxyMode int

	// 1.直接寻址： url = upstream_value
	// 2.间接寻址: 在etcd中取出key为Api.Name的值，返回的数据结构存储在UpstreamValue
	UpstreamMode    int
	UpstreamServers []*UpstreamServer
}

type UpstreamServer struct {
	IP   string
	Load int
}

var Apis = &sync.Map{}
