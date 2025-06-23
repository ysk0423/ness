package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTradeNote(t *testing.T) {
	tradeDate := time.Now()
	symbol := "AAPL"
	content := "購入理由: 業績好調"

	tradeNote := NewTradeNote(tradeDate, symbol, content)

	assert.NotNil(t, tradeNote)
	assert.Equal(t, tradeDate, tradeNote.TradeDate)
	assert.Equal(t, symbol, tradeNote.Symbol)
	assert.Equal(t, content, tradeNote.Content)
	assert.Equal(t, uint(0), tradeNote.ID) // 新規作成時はID=0
}

func TestTradeNote_Update(t *testing.T) {
	// 初期データを設定
	originalDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tradeNote := NewTradeNote(originalDate, "AAPL", "購入理由: 業績好調")

	// 更新データを設定
	newDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	newSymbol := "GOOGL"
	newContent := "売却理由: 利確"

	// 更新を実行
	tradeNote.Update(newDate, newSymbol, newContent)

	// 更新後の値を検証
	assert.Equal(t, newDate, tradeNote.TradeDate)
	assert.Equal(t, newSymbol, tradeNote.Symbol)
	assert.Equal(t, newContent, tradeNote.Content)
}

func TestTradeNote_Fields(t *testing.T) {
	tradeDate := time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC)
	symbol := "TSLA"
	content := "技術的分析に基づく購入"

	tradeNote := NewTradeNote(tradeDate, symbol, content)

	// フィールドの型と値を確認
	assert.IsType(t, uint(0), tradeNote.ID)
	assert.IsType(t, time.Time{}, tradeNote.TradeDate)
	assert.IsType(t, "", tradeNote.Symbol)
	assert.IsType(t, "", tradeNote.Content)
	assert.IsType(t, time.Time{}, tradeNote.CreatedAt)
	assert.IsType(t, time.Time{}, tradeNote.UpdatedAt)

	// 値の確認
	assert.Equal(t, tradeDate, tradeNote.TradeDate)
	assert.Equal(t, symbol, tradeNote.Symbol)
	assert.Equal(t, content, tradeNote.Content)
}