package portal

import "net/http"

const (
	PortalServiceHost = "127.0.0.1"
	PortalServicePort = 5000
	RemoteUpdateUrl   = "remote-update"
)

func RegisterHandlers() {
	http.HandleFunc(RemoteUpdateUrl, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// TODO 更新协议，服务端需要确定
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
