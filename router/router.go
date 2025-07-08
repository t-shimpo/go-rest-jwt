package router

import (
	"net/http"

	"github.com/t-shimpo/go-rest-jwt/auth"
	"github.com/t-shimpo/go-rest-jwt/handlers"
)

func methodNotAllowedHandler(w http.ResponseWriter) {
	http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed)
}

func SetupRoutes(userHandler *handlers.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		case http.MethodGet:
			userHandler.GetUsers(w, r)
		default:
			methodNotAllowedHandler(w)
		}
	})

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByID(w, r)
		case http.MethodPatch:
			userHandler.PatchUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		default:
			methodNotAllowedHandler(w)
		}
	})

	mux.Handle("/users/", auth.JTWMiddleware(protectedHandler))

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userHandler.Login(w, r)
		default:
			methodNotAllowedHandler(w)
		}
	})

	return mux
}
