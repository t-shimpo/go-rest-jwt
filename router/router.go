package router

import (
	"net/http"

	"github.com/t-shimpo/go-rest-standard-library/handlers"
)

func methodNotAllowedHandler(w http.ResponseWriter) {
	http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed)
}

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateUserHandler(w, r)
		case http.MethodGet:
			handlers.GetUsersHandler(w, r)
		default:
			methodNotAllowedHandler(w)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetUserHandler(w, r)
		default:
			methodNotAllowedHandler(w)
		}
	})

	return mux
}
