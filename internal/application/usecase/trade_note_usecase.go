package usecase

import (
	"context"

	"ness/internal/application/dto"
	"ness/internal/domain/entity"
	"ness/internal/domain/repository"
)

type TradeNoteUsecase struct {
	tradeNoteRepo repository.TradeNoteRepository
}

func NewTradeNoteUsecase(tradeNoteRepo repository.TradeNoteRepository) *TradeNoteUsecase {
	return &TradeNoteUsecase{
		tradeNoteRepo: tradeNoteRepo,
	}
}

func (u *TradeNoteUsecase) CreateTradeNote(ctx context.Context, req *dto.TradeNoteRequest) (*dto.TradeNoteResponse, error) {
	tradeNote := entity.NewTradeNote(req.TradeDate, req.Symbol, req.Content)
	
	if err := u.tradeNoteRepo.Create(ctx, tradeNote); err != nil {
		return nil, err
	}

	return &dto.TradeNoteResponse{
		ID:        tradeNote.ID,
		TradeDate: tradeNote.TradeDate,
		Symbol:    tradeNote.Symbol,
		Content:   tradeNote.Content,
		CreatedAt: tradeNote.CreatedAt,
		UpdatedAt: tradeNote.UpdatedAt,
	}, nil
}

func (u *TradeNoteUsecase) GetTradeNote(ctx context.Context, id uint) (*dto.TradeNoteResponse, error) {
	tradeNote, err := u.tradeNoteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.TradeNoteResponse{
		ID:        tradeNote.ID,
		TradeDate: tradeNote.TradeDate,
		Symbol:    tradeNote.Symbol,
		Content:   tradeNote.Content,
		CreatedAt: tradeNote.CreatedAt,
		UpdatedAt: tradeNote.UpdatedAt,
	}, nil
}

func (u *TradeNoteUsecase) GetAllTradeNotes(ctx context.Context) (*dto.TradeNoteListResponse, error) {
	tradeNotes, err := u.tradeNoteRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.TradeNoteResponse
	for _, tradeNote := range tradeNotes {
		responses = append(responses, dto.TradeNoteResponse{
			ID:        tradeNote.ID,
			TradeDate: tradeNote.TradeDate,
			Symbol:    tradeNote.Symbol,
			Content:   tradeNote.Content,
			CreatedAt: tradeNote.CreatedAt,
			UpdatedAt: tradeNote.UpdatedAt,
		})
	}

	return &dto.TradeNoteListResponse{
		TradeNotes: responses,
		Total:      len(responses),
	}, nil
}

func (u *TradeNoteUsecase) UpdateTradeNote(ctx context.Context, id uint, req *dto.TradeNoteRequest) (*dto.TradeNoteResponse, error) {
	tradeNote, err := u.tradeNoteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tradeNote.Update(req.TradeDate, req.Symbol, req.Content)

	if err := u.tradeNoteRepo.Update(ctx, tradeNote); err != nil {
		return nil, err
	}

	return &dto.TradeNoteResponse{
		ID:        tradeNote.ID,
		TradeDate: tradeNote.TradeDate,
		Symbol:    tradeNote.Symbol,
		Content:   tradeNote.Content,
		CreatedAt: tradeNote.CreatedAt,
		UpdatedAt: tradeNote.UpdatedAt,
	}, nil
}

func (u *TradeNoteUsecase) DeleteTradeNote(ctx context.Context, id uint) error {
	return u.tradeNoteRepo.Delete(ctx, id)
}