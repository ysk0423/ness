package main

import (
	"ness/internal/presentation/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ハンドラーを初期化
	h := handler.NewHelloHandler()

	// ルートを設定
	e.GET("/", h.Hello)

	// サーバーを開始
	e.Logger.Fatal(e.Start(":8080"))
}