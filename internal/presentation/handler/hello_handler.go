package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloHandler Hello World関連のハンドラー
type HelloHandler struct{}

// NewHelloHandler HelloHandlerのコンストラクタ
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// Hello Hello Worldメッセージを返す
func (h *HelloHandler) Hello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Hello, World!",
	})
}
