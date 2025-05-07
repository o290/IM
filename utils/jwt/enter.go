package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtPayload struct {
	UserID   uint   `json:"userID"`
	NickName string `json:"nickName"`
	Role     int8   `json:"role"`
}
type CustomClaims struct {
	JwtPayload
	jwt.RegisteredClaims
}

// GenToken 生成token
func GenToken(payload JwtPayload, accessSecret string, expires int64) (string, error) {
	//1.创建自定义的结构体
	claim := CustomClaims{
		JwtPayload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expires))),
		},
	}
	//2.创建jwt对象,指明使用的签名算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	//3.SignedString进行签名生成签名字符串
	return token.SignedString([]byte(accessSecret))
}

// ParseToken 验证token
func ParseToken(tokenStr string, accessSecret string) (*CustomClaims, error) {
	//1.根据私钥对jwt进行解析
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	//2.类型断言并判断令牌是否有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
