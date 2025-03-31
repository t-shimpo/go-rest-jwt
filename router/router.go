package router

import (
	"net/http"

	"github.com/t-shimpo/go-rest-standard-library/handlers"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.CreateUserHandler(w, r)
		} else {
			http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
