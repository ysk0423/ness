package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"ness/internal/application/usecase"
	"ness/internal/infrastructure/database"
	"ness/internal/infrastructure/repository"
	"ness/internal/presentation/handler"
)

func main() {
	// データベース接続を初期化
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// リポジトリを初期化
	tradeNoteRepo := repository.NewTradeNoteRepository(db.DB)

	// ユースケースを初期化
	tradeNoteUsecase := usecase.NewTradeNoteUsecase(tradeNoteRepo)

	// ハンドラーを初期化
	helloHandler := handler.NewHelloHandler()
	tradeNoteHandler := handler.NewTradeNoteHandler(tradeNoteUsecase)

	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ルートを設定
	e.GET("/", helloHandler.Hello)

	// トレードノートAPIルート
	e.GET("/notes", tradeNoteHandler.GetTradeNotes)
	e.POST("/notes", tradeNoteHandler.CreateTradeNote)
	e.GET("/notes/:id", tradeNoteHandler.GetTradeNote)
	e.PUT("/notes/:id", tradeNoteHandler.UpdateTradeNote)
	e.DELETE("/notes/:id", tradeNoteHandler.DeleteTradeNote)

	// サーバーを開始
	e.Logger.Fatal(e.Start(":8080"))
}
