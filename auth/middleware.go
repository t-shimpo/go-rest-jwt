package auth

import (
	"net/http"
	"strings"
)

func JTWMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーからトークンを取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// subject (user ID) をリクエストのコンテキストに設定
		userID, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// ContextにユーザーIDを設定
		ctx := SetUserID(r.Context(), int(userID))

		// 次のハンドラーを呼び出す
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
