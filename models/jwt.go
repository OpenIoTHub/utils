package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type TokenClaims struct {
	RunId      string
	Host       string
	TcpPort    int
	KcpPort    int
	TlsPort    int
	GrpcPort   int
	P2PApiPort int
	Permission int
	jwt.StandardClaims
}

func GetToken(gatewayConfig GatewayConfig, permission int, expiresecd int64) (token string, err error) {
	tokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		gatewayConfig.LastId,
		gatewayConfig.Server.ServerHost,
		gatewayConfig.Server.TcpPort,
		gatewayConfig.Server.KcpPort,
		gatewayConfig.Server.TlsPort,
		gatewayConfig.Server.GrpcPort,
		gatewayConfig.Server.UdpApiPort,
		permission,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 8*60*60,
			ExpiresAt: time.Now().Unix() + expiresecd,
		},
	})
	tokenStr, err := tokenModel.SignedString([]byte(gatewayConfig.Server.LoginKey))
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func DecodeToken(salt, tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(salt), nil
	})
	if err != nil {
		log.Println("错误")
		return &TokenClaims{}, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		//log.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		return &TokenClaims{}, fmt.Errorf("jwt decode err")
	}
}

func DecodeUnverifiedToken(tokenStr string) (*TokenClaims, error) {
	token, _ := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(""), nil
	})
	//不校验是否是正确加密的signature
	if token == nil || token.Claims == nil {
		return &TokenClaims{}, fmt.Errorf("token or token.Claims is nil")
	}
	if claims, ok := token.Claims.(*TokenClaims); ok {
		//log.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		return &TokenClaims{}, fmt.Errorf("jwt decode err，not TokenClaims")
	}
}
