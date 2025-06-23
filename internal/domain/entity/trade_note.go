package entity

import (
	"time"

	"gorm.io/gorm"
)

type TradeNote struct {
	ID          uint           `gorm:"primarykey"`
	TradeDate   time.Time      `gorm:"not null"`
	Symbol      string         `gorm:"not null;size:50"`
	Content     string         `gorm:"not null;type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func NewTradeNote(tradeDate time.Time, symbol string, content string) *TradeNote {
	return &TradeNote{
		TradeDate: tradeDate,
		Symbol:    symbol,
		Content:   content,
	}
}

func (tn *TradeNote) Update(tradeDate time.Time, symbol string, content string) {
	tn.TradeDate = tradeDate
	tn.Symbol = symbol
	tn.Content = content
}