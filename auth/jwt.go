package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/t-shimpo/go-rest-jwt/config"
)

// 受け取ったユーザーIDからJWTトークンを生成する関数
func GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                                                                      // subjectとしてユーザーIDを設定
		"iss": config.JWTIssuer,                                                            // 発行者
		"exp": time.Now().Add(time.Duration(config.JWTExpireMinutes) * time.Minute).Unix(), // 有効期限
		"iat": time.Now().Unix(),                                                           // 発行日時
	}

	// HS256で署名するトークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵を使ってトークンに署名
	return token.SignedString([]byte(config.JWTSecret))
}
