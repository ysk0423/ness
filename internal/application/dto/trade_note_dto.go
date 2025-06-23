package dto

import "time"

type TradeNoteRequest struct {
	TradeDate time.Time `json:"trade_date" validate:"required"`
	Symbol    string    `json:"symbol" validate:"required,max=50"`
	Content   string    `json:"content" validate:"required"`
}

type TradeNoteResponse struct {
	ID        uint      `json:"id"`
	TradeDate time.Time `json:"trade_date"`
	Symbol    string    `json:"symbol"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TradeNoteListResponse struct {
	TradeNotes []TradeNoteResponse `json:"trade_notes"`
	Total      int                 `json:"total"`
}