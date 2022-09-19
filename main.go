package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	// Claimsオブジェクトの作成
	claims := jwt.MapClaims{
		"user_id": 12345678,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// ヘッダーとペイロードの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("Header: %#v\n", token.Header) // Header: map[string]interface {}{"alg":"HS256", "typ":"JWT"}
	fmt.Printf("Claims: %#v\n", token.Claims) // CClaims: jwt.MapClaims{"exp":1634051243, "user_id":12345678}

	// トークンに署名を付与
	tokenString, _ := token.SignedString([]byte("SECRET_KEY"))
	fmt.Println("tokenString:", tokenString) // tokenString: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzQwNTEzNjMsInVzZXJfaWQiOjEyMzQ1Njc4fQ.OooYrharapD5X2LV5UUWBOkEqH57wDfMd5ibkIpJHYM
}
