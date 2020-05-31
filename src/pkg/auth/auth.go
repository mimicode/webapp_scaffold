package auth

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(val string) string {
	h := md5.New()
	h.Write([]byte(val)) // 先加盐
	return hex.EncodeToString(h.Sum([]byte(val)))
}
