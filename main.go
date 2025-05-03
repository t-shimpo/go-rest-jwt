package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/t-shimpo/go-rest-standard-library-layered/config"
	"github.com/t-shimpo/go-rest-standard-library-layered/handlers"
	"github.com/t-shimpo/go-rest-standard-library-layered/repository"
	"github.com/t-shimpo/go-rest-standard-library-layered/service"

	"github.com/t-shimpo/go-rest-standard-library-layered/router"
)

func main() {
	// DB 初期化
	db, err := config.InitDB()
	if err != nil {
		fmt.Println("DB初期化エラー:", err)
		return
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	mux := router.SetupRoutes(userHandler)

	fmt.Println("サーバーが 8080 ポートで起動中")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
