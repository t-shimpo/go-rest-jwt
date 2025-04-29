package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/t-shimpo/go-rest-standard-library-layered/config"

	"github.com/t-shimpo/go-rest-standard-library-layered/router"
)

func main() {
	// DB 初期化
	err := config.InitDB()
	if err != nil {
		fmt.Println("DB初期化エラー:", err)
		return
	}

	if config.DB != nil {
		defer func() {
			if err := config.DB.Close(); err != nil {
				fmt.Println("DBクローズエラー:", err)
			}
		}()
	}

	// ルーティング設定
	mux := router.SetupRoutes()

	// サーバー起動
	fmt.Println("サーバーが 8080 ポートで起動中")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
