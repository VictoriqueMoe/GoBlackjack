package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

var h = sha256.New()

func DeviceHash(ip, ua string) string {
	h := sha256.New()
	h.Write([]byte(ip + ua))
	return hex.EncodeToString(h.Sum(nil))
}

func Pop[T any](list *[]T) T {
	f := len(*list)
	rv := (*list)[f-1]
	*list = (*list)[:f-1]
	return rv
}
