package models

// 容器的网络资源
type IpConfig struct {
	IP   string `json:"IP"`   // IP
	Mask string `json:"Mask"` // 掩码
}
type ContainerNetWorkConfig struct {
	Name    string   `json:"name"` // 网络名
	IpRange IpConfig `json:"IpRange"`
	Driver  string   `json:"Driver"` // 驱动
}
