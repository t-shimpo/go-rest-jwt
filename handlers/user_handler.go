package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/t-shimpo/go-rest-jwt/auth"
	"github.com/t-shimpo/go-rest-jwt/models"
	"github.com/t-shimpo/go-rest-jwt/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "無効なリクエストボディ")
		return
	}

	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	createdUser, err := h.userService.CreateUser(&user, req.Password)
	if err != nil {
		if err == service.ErrValidation {
			respondWithError(w, http.StatusBadRequest, "入力値が不正です")
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "ユーザー作成に失敗しました")
			return
		}
	}

	respondWithJson(w, http.StatusCreated, createdUser)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// URLパスからIDを取得
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "IDは必要です")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "IDは数値である必要があります")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ユーザー取得時にエラーが発生しました")
		return
	}
	if user == nil {
		respondWithError(w, http.StatusNotFound, "ユーザーが見つかりません")
		return
	}

	respondWithJson(w, http.StatusOK, user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // デフォルト値
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // デフォルト値
	}

	users, err := h.userService.GetUsers(limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ユーザー一覧の取得に失敗しました")
		return
	}

	respondWithJson(w, http.StatusOK, users)
}

func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	// URLパスからIDを取得
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "IDは必要です")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "IDは数値である必要があります")
		return
	}

	// リクエストボディのパース
	var req struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	if req.Name == nil && req.Email == nil {
		respondWithError(w, http.StatusBadRequest, "更新するフィールドを指定してください")
		return
	}

	updatedUser, err := h.userService.PatchUser(id, req.Name, req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "ユーザー更新中にエラーが発生しました")
		return
	}
	if updatedUser == nil {
		respondWithError(w, http.StatusNotFound, "ユーザーが見つかりません")
		return
	}

	respondWithJson(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// URLパスからIDを取得
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "IDは必要です")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "IDは数値である必要があります")
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "ユーザーが見つかりません")
		} else {
			respondWithError(w, http.StatusInternalServerError, "ユーザー削除中にエラーが発生しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "無効なリクエストボディ")
		return
	}

	user, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		if err == service.ErrorNotFound || err == service.ErrorInvalidPassword {
			respondWithError(w, http.StatusUnauthorized, "認証に失敗しました")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "ログイン処理に失敗しました")
		return
	}

	token, err := auth.GenerateToken(int64(user.ID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "トークン生成に失敗しました")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func respondWithJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJson(w, status, map[string]string{"error": message})
}
