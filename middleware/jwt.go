package middleware

import (
	"Team2048_Tiktok/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//用于加密签名的密钥
var myKey = []byte("Tiktok2048")

type MyClaims struct {
	UserId int64
	jwt.StandardClaims
}

// ReleaseToken	颁发JWT token
func ReleaseToken(user model.User) (string, error) {
	//创建一个自己的声明
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &MyClaims{
		UserId: user.Id,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Tiktok2048",
		}}

	// 使用指定的签名方法创建签名对象，一般使用HS256加密
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的密钥进行签名，并获得完整的编码后的字符串token
	return token.SignedString(myKey)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*MyClaims, error) {
	var claims = new(MyClaims)

	// 解析token字符串，获取jwt.Token 以及 信息结构体mc
	//提供服务端的密钥进行解密
	decodeFunc := func(token *jwt.Token) (i interface{}, err error) {
		return myKey, nil
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, decodeFunc)
	if err != nil {
		return nil, err
	}

	// 校验token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}

}
