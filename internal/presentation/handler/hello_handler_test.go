package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHelloHandler_Hello(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "正常なHello Worldレスポンス",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Hello, World!"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoインスタンスを作成
			e := echo.New()

			// リクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// ハンドラーを作成
			handler := NewHelloHandler()

			// テスト実行
			err := handler.Hello(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
