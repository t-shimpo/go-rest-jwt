package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/t-shimpo/go-rest-standard-library/models"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// エラーレスポンス
func respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// `POST /users“ の処理
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	var req CreateUserRequest
	defer r.Body.Close()

	// json デコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "無効なリクエストボディ")
		return
	}

	// バリデーション
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "名前は必須です")
		return
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "メールは必須です")
		return
	}

	// DB に保存
	user, err := models.CreateUser(req.Name, req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ユーザー作成中にエラーが発生しました")
		return
	}

	// 成功レスポンス
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
