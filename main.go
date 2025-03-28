package main

import (
	"fmt"

	"github.com/t-shimpo/go-rest-standard-library/config"
)

func main() {
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
}
