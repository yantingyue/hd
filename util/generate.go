package util

import (
	"encoding/base64"
	"time"

	"demo/hd/constant"
	"github.com/spf13/cast"
)

func GenerateHeader(token string) map[string]string {
	timestamp := cast.ToString(time.Now().UnixMilli())
	return map[string]string{
		"token":     token,
		"version":   constant.Version,
		"channel":   constant.Channel,
		"platform":  constant.Platform,
		"appname":   constant.Appname,
		"timestamp": timestamp,
		"sign":      MD5(timestamp + constant.Salt),
	}
}

func GenerateCreateOrderHeader(token string) map[string]string {
	timestamp := cast.ToString(time.Now().UnixMilli())
	sessionBytes, _ := Aes.Encrypt([]byte(timestamp))
	return map[string]string{
		"token":     token,
		"version":   constant.Version,
		"channel":   constant.Channel,
		"platform":  constant.Platform,
		"appname":   constant.Appname,
		"timestamp": timestamp,
		"sign":      MD5(timestamp + constant.Salt),
		"session":   base64.StdEncoding.EncodeToString(sessionBytes),
		"cookie":    CreateOrderCookie(timestamp, token, base64.StdEncoding.EncodeToString(sessionBytes)),
	}
}
