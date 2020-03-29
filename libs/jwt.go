package libs

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtCustomClaims struct {
	jwt.StandardClaims

	// 追加自己需要的信息
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

/**
生成 token
SecretKey 是一个 const 常量
*/
func CreateToken(SecretKey []byte, id uint, name, email string) (tokenString string, err error) {
	claims := &JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "chenishr@gmail.com",
			Subject:   "Cycron",
		},
		Id:       id,
		UserName: name,
		Email:    email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

/**
解析 token
*/
func ParseToken(tokenSrt string, SecretKey []byte) (claims *JwtCustomClaims, err error) {
	var token *jwt.Token

	token, err = jwt.ParseWithClaims(tokenSrt, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims = token.Claims.(*JwtCustomClaims)
	return
}
