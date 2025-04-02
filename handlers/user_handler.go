package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// `POST /users`
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

// `GET /users/{id}`
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// URL から ID を取得
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 {
		respondWithError(w, http.StatusBadRequest, "無効なURLです")
		return
	}
	idStr := segments[2]

	// ID を整数に変換
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "IDは数値である必要があります")
		return
	}

	// DB からユーザー取得
	user, err := models.GetUserByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "ユーザーが見つかりません")
		} else {
			respondWithError(w, http.StatusInternalServerError, "ユーザー取得中にエラーが発生しました")
		}
		return
	}

	// 成功レスポンス
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
