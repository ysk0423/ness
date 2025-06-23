package repository

import (
	"context"
	"ness/internal/domain/entity"
)

type TradeNoteRepository interface {
	Create(ctx context.Context, tradeNote *entity.TradeNote) error
	GetByID(ctx context.Context, id uint) (*entity.TradeNote, error)
	GetAll(ctx context.Context) ([]*entity.TradeNote, error)
	Update(ctx context.Context, tradeNote *entity.TradeNote) error
	Delete(ctx context.Context, id uint) error
}