package registry

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"sync"
	"time"

	"github.com/draymonders/distributed/util"
)

const (
	RegisterServiceHost = "localhost"
	RegisterServicePort = ":3000"
	RegisterServicePath = "/services"
)

// Registration 服务端信息
type Registration struct {
	ServiceMap map[ServiceName][]*Service
	mux        *sync.RWMutex
}

var (
	registration Registration
)

func InitRegistry() {
	registration.ServiceMap = make(map[ServiceName][]*Service, 0)
	registration.mux = new(sync.RWMutex)
	go registration.HeartBeat(2 * time.Second)
}

// RegisterHandlers 实现http 方法
func RegisterHandlers() {
	http.HandleFunc(RegisterServicePath, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		// 新增服务
		case http.MethodPost:
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			service := Service{}
			err = json.Unmarshal(data, &service)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = registration.add(&service)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			registerInfo := service.transToRegisterInfo(true)
			go registration.NotifyAll(registerInfo)
			stdlog.Printf("service [%v] registerd.", service.Name)
			return
		// 获取所有服务
		case http.MethodGet:
			serviceList := registration.toList()
			data, err := json.Marshal(serviceList)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
			return
		// 删除服务
		case http.MethodDelete:
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			service := Service{}
			err = json.Unmarshal(data, &service)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = registration.del(service)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			stdlog.Printf("service [%v] be deleted.", service.Name)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func (r *Registration) add(service *Service) error {
	if service == nil {
		return nil
	}
	r.mux.Lock()
	defer r.mux.Unlock()

	name := service.Name
	if _, ok := r.ServiceMap[name]; !ok {
		r.ServiceMap[name] = make([]*Service, 0)
	}
	r.ServiceMap[name] = append(r.ServiceMap[name], service)
	// notify()
	return nil
}

func (r *Registration) del(service Service) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	name := service.Name
	if _, ok := r.ServiceMap[name]; !ok {
		return nil
	}
	for i, svc := range r.ServiceMap[name] {
		if svc.Url == service.Url {
			r.ServiceMap[name] = append(r.ServiceMap[name][:i], r.ServiceMap[name][i+1:]...)
		}
	}
	return nil
}

func (r *Registration) toList() []*Service {
	services := make([]*Service, 0)
	for _, serviceList := range r.ServiceMap {
		services = append(services, serviceList...)
	}
	return services
}

// NotifyAll 广播通知客户端，服务新发现/注销，这应该有个单独的结构体去封装
func (r *Registration) NotifyAll(info *RegisterInfo) {
	if info == nil {
		return
	}
	sendToClientInfos := info.Added
	if len(info.Removed) != 0 {
		sendToClientInfos = append(sendToClientInfos, info.Removed...)
	}
	if len(sendToClientInfos) == 0 {
		return
	}
	// 只需要向不为同名的服务发送更新信息
	for _, curSvc := range sendToClientInfos {
		curSvcName := curSvc.Name
		for svcName, svcs := range r.ServiceMap {
			if svcName == curSvcName {
				continue
			}
			for _, svc := range svcs {
				if svc.UpdateUrl != "" {
					notify(curSvc, svc.Name, svc.UpdateUrl)
				}
			}
		}
	}
}

func notify(info *ClientServiceInfo, svcName ServiceName, clientUpdateUrl string) {
	if info == nil {
		return
	}
	svcInfoBytes, err := json.Marshal(info)
	if err != nil {
		stdlog.Printf("json.Marshal clientServiceInfo fail, err: %v", err)
		return
	}
	resp, err := http.Post(clientUpdateUrl, "application/json", bytes.NewBuffer(svcInfoBytes))
	if err != nil {
		stdlog.Printf("notify to client [%v] fail, err: %v", svcName, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		stdlog.Printf("notify to client [%v] success", svcName)
	} else {
		stdlog.Printf("notify to client [%v] fail, resp: %v", svcName, util.MarshalJsonNotErr(resp.Body))
	}
}

// HeartBeat 心跳服务，检测所注册的服务的状态
func (r *Registration) HeartBeat(d time.Duration) {
	const retryCnt = 3
	const retryTimeDuration = 100 * time.Millisecond
	for {
		// 遍历所有的svcName
		for svcName, svcs := range r.ServiceMap {
			for i, svc := range svcs {
				// 重试心跳，如若仍失败，则从注册中心中剔除该实例
				isAlive := util.DoRetry(retryCnt, retryTimeDuration, func() bool {
					return doHeartBeat(svcName, svc)
				})
				if !isAlive {
					r.mux.Lock()
					r.ServiceMap[svcName] = append(r.ServiceMap[svcName][:i], r.ServiceMap[svcName][i+1:]...)
					r.mux.Unlock()
				}
			}
		}
		time.Sleep(d)
	}
}

func doHeartBeat(name ServiceName, svc *Service) bool {
	if svc == nil || svc.HeartBeatUrl == "" {
		return false
	}
	resp, err := http.Get(svc.HeartBeatUrl)
	if err != nil {
		stdlog.Printf("service [%v] ping failed, err: %v", name, err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		stdlog.Printf("service [%v] ping success", name)
		return true
	}
	return false
}
