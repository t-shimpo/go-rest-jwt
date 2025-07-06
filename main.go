package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/t-shimpo/go-rest-jwt/config"
	"github.com/t-shimpo/go-rest-jwt/handlers"
	"github.com/t-shimpo/go-rest-jwt/repository"
	"github.com/t-shimpo/go-rest-jwt/service"

	"github.com/t-shimpo/go-rest-jwt/router"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("環境変数の読み込みに失敗しました:", err)
	}
}

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
