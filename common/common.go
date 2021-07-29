package common

import "net/http"

func HeartBeatFunc(w http.ResponseWriter, r *http.Request) {
	content := "pong"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
