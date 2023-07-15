package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token is not active yet")
	TokenMalformed   = errors.New("token is malformed")
	TokenInvalid     = errors.New("token is invalid")
)

const (
	SigningKey  = "gameapi2qwessdjsghegj"
	ExpiresTime = 604800 // 7days
	//ExpiresTime = 3600 // 1h
	// BufferTime  = 86400 // 24h
	BufferTime = 600
	Issuer     = "game-api"
)

//func main() {
//	fmt.Println(GetToken(160487))
//	//fmt.Println(GetToken(160534))
//	//fmt.Println(GetToken(160535))
//	//t1("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIxNjAzODEiLCJCdWZmZXJUaW1lIjo2MDAsImV4cCI6MTY3OTI4Mjk3NSwiaXNzIjoiZ2FtZS1hcGkiLCJuYmYiOjE2Nzg2NzcxNzV9.AlLJa1wEdsMgRYcsfPaOrocUqJehdGaKKuU7pkyVw74")
//	//t1("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiIxNjAzNzQiLCJCdWZmZXJUaW1lIjo2MDAsImV4cCI6MTY4MTA5MTYxOCwiaXNzIjoiZ2FtZS1hcGkiLCJuYmYiOjE2ODA0ODU4MTh9.ncOnb8nHNFQGA1k1mkjkleJYPzKx9-nxvpTBzA7i8hY")
//}

func t1(token string) {
	// 解析token包含的信息
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		return
	}
	fmt.Println(claims)
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(SigningKey),
	}
}

type BaseClaims struct {
	UserId string
}

// Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,        // 签名生效时间
			ExpiresAt: time.Now().Unix() + ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    Issuer,                          // 签名的发行者
		},
	}
	return claims
}

func GetToken(userId uint64) (token string, err error) {
	j := JWT{SigningKey: []byte(SigningKey)} // 唯一签名
	claims := j.CreateClaims(BaseClaims{
		UserId: strconv.FormatUint(userId, 10),
	})
	token, err = j.CreateToken(claims)
	return
}

// 创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
