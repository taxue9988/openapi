package apidata

/* 从Mysql中加载API信息到内存中*/
type API struct {
	ID int `json:"id"`

	FullName  string `json:"api_name"`
	Company   string `json:"company"`
	Product   string `json:"product"`
	System    string `json:"system"`
	Interface string `json:"interface"`
	Version   string `json:"version"`

	Method string `json:"method"`

	ProxyMode string `json:"proxy_mode"`

	UpstreamMode  string `json:"upstream_mode"`
	UpstreamValue string `json:"upstream_value"`
}
