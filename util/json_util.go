package util

import "encoding/json"

func MarshalJsonNotErr(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
