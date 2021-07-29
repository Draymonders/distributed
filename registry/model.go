package registry

import "fmt"

type ServiceName string

type Service struct {
	Name               ServiceName   `json:"name"`
	Url                string        `json:"url"`                 // 服务对应的url
	DependencyServices []ServiceName `json:"dependency_services"` // 依赖的服务名
	UpdateUrl          string        `json:"update_url"`          // 被注册服务需要接受其他服务更新的url
	HeartBeatUrl       string        `json:"heart_beat_url"`      // 被注册服务的心跳url
	Methods            []*Method     `json:"methods"`             // 服务拥有的方法
}

type Method struct {
	Name   string `json:"name"`   // 方法名，匹配到服务名后，接着匹配到方法名，即可打入请求
	Path   string `json:"path"`   // 方法对应的url后缀
	Weight int64  `json:"weight"` // 方法权重，冷启动时值较小
}

// ClientServiceInfo 客户端服务信息
type ClientServiceInfo struct {
	Name    ServiceName `json:"name"`
	Urls    []string    `json:"urls"`
	Methods []*Method   `json:"methods"`
}

type RegisterInfo struct {
	Added   []*ClientServiceInfo
	Removed []*ClientServiceInfo
}

func NewService(name, host string, port int) *Service {
	return &Service{
		Name:         ServiceName(name),
		Url:          fmt.Sprintf("http://%s:%d", host, port),
		HeartBeatUrl: fmt.Sprintf("http://%s:%d/ping", host, port),
		UpdateUrl:    fmt.Sprintf("http://%s:%d/remote-update", host, port),
	}
}

func (s *Service) transToRegisterInfo(isAdd bool) *RegisterInfo {
	registerInfo := &RegisterInfo{}
	clientServiceInfo := &ClientServiceInfo{
		Name:    s.Name,
		Urls:    []string{s.Url},
		Methods: s.Methods,
	}
	clientServiceInfos := []*ClientServiceInfo{
		clientServiceInfo,
	}

	if isAdd {
		registerInfo.Added = clientServiceInfos
	} else {
		registerInfo.Removed = clientServiceInfos
	}
	return registerInfo
}
