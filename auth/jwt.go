package auth

import (
	"errors"
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

// JWTトークンを検証し、Claimsを返す関数
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 署名方式の確認
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// 秘密鍵を返す
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// ClaimsをMapClaimsにキャストし、トークンが有効かどうかを確認
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 有効期限を検証
	if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
		return nil, errors.New("token has expired")
	}

	// 発行者を検証
	if iss, ok := claims["iss"].(string); !ok || iss != config.JWTIssuer {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}
