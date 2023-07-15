package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"github.com/speps/go-hashids/v2"
	"golang.org/x/crypto/bcrypt"
)

const (
	HashSalt = "7cd30abtzna8"
)

func HashPassword(pass *string) {
	bytePass := []byte(*pass)
	hPass, _ := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	*pass = string(hPass)
}

func ComparePassword(dbPass, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass)) == nil
}

func Encode(id int) (e string) {
	hd := hashids.NewData()
	hd.Salt = HashSalt
	h, err := hashids.NewWithData(hd)
	if nil != err {
		return ""
	}
	e, err = h.Encode([]int{id})
	if nil != err {
		return ""
	}
	return e
}

func EncodeInt64(id int64) (e string) {
	hd := hashids.NewData()
	hd.Salt = HashSalt
	h, err := hashids.NewWithData(hd)
	if nil != err {
		return ""
	}
	e, err = h.EncodeInt64([]int64{id})
	if nil != err {
		return ""
	}
	return e
}

func Decode(id string) (d int) {
	hd := hashids.NewData()
	hd.Salt = HashSalt
	h, err := hashids.NewWithData(hd)
	if nil != err {
		return 0
	}
	var ds []int
	ds, err = h.DecodeWithError(id)
	if nil != err {
		return 0
	}
	if len(ds) == 0 {
		return 0
	}
	return ds[0]
}

func DecodeInt64(id string) (d int64) {
	hd := hashids.NewData()
	hd.Salt = HashSalt
	h, err := hashids.NewWithData(hd)
	if nil != err {
		return 0
	}
	var ds []int64
	ds, err = h.DecodeInt64WithError(id)
	if nil != err {
		return 0
	}
	if len(ds) == 0 {
		return 0
	}
	return ds[0]
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func CreateOrderCookie(timeStamp, token, sessionId string) string {
	b, _ := base64.StdEncoding.DecodeString(sessionId)
	a, _ := Aes.Decrypt(b)
	s := "286ee9871c1e4f4682d27722a3c51c83" + timeStamp + token + string(a)
	return MD5(s)
}
