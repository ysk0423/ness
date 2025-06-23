package repository

import (
	"context"

	"gorm.io/gorm"

	"ness/internal/domain/entity"
	"ness/internal/domain/repository"
)

type tradeNoteRepositoryImpl struct {
	db *gorm.DB
}

func NewTradeNoteRepository(db *gorm.DB) repository.TradeNoteRepository {
	return &tradeNoteRepositoryImpl{
		db: db,
	}
}

func (r *tradeNoteRepositoryImpl) Create(ctx context.Context, tradeNote *entity.TradeNote) error {
	return r.db.WithContext(ctx).Create(tradeNote).Error
}

func (r *tradeNoteRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.TradeNote, error) {
	var tradeNote entity.TradeNote
	err := r.db.WithContext(ctx).First(&tradeNote, id).Error
	if err != nil {
		return nil, err
	}
	return &tradeNote, nil
}

func (r *tradeNoteRepositoryImpl) GetAll(ctx context.Context) ([]*entity.TradeNote, error) {
	var tradeNotes []*entity.TradeNote
	err := r.db.WithContext(ctx).Order("trade_date desc").Find(&tradeNotes).Error
	if err != nil {
		return nil, err
	}
	return tradeNotes, nil
}

func (r *tradeNoteRepositoryImpl) Update(ctx context.Context, tradeNote *entity.TradeNote) error {
	return r.db.WithContext(ctx).Save(tradeNote).Error
}

func (r *tradeNoteRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.TradeNote{}, id).Error
}