package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ness/internal/application/dto"
	"ness/internal/domain/entity"
)

// MockTradeNoteRepository is a mock implementation of TradeNoteRepository
type MockTradeNoteRepository struct {
	mock.Mock
}

func (m *MockTradeNoteRepository) Create(ctx context.Context, tradeNote *entity.TradeNote) error {
	args := m.Called(ctx, tradeNote)
	return args.Error(0)
}

func (m *MockTradeNoteRepository) GetByID(ctx context.Context, id uint) (*entity.TradeNote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TradeNote), args.Error(1)
}

func (m *MockTradeNoteRepository) GetAll(ctx context.Context) ([]*entity.TradeNote, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.TradeNote), args.Error(1)
}

func (m *MockTradeNoteRepository) Update(ctx context.Context, tradeNote *entity.TradeNote) error {
	args := m.Called(ctx, tradeNote)
	return args.Error(0)
}

func (m *MockTradeNoteRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTradeNoteUsecase_CreateTradeNote(t *testing.T) {
	mockRepo := new(MockTradeNoteRepository)
	usecase := NewTradeNoteUsecase(mockRepo)
	ctx := context.Background()

	req := &dto.TradeNoteRequest{
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "購入理由: 業績好調",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.TradeNote")).Return(nil).Run(func(args mock.Arguments) {
		tradeNote := args.Get(1).(*entity.TradeNote)
		tradeNote.ID = 1
		tradeNote.CreatedAt = time.Now()
		tradeNote.UpdatedAt = time.Now()
	})

	result, err := usecase.CreateTradeNote(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, req.TradeDate, result.TradeDate)
	assert.Equal(t, req.Symbol, result.Symbol)
	assert.Equal(t, req.Content, result.Content)
	mockRepo.AssertExpectations(t)
}

func TestTradeNoteUsecase_CreateTradeNote_Error(t *testing.T) {
	mockRepo := new(MockTradeNoteRepository)
	usecase := NewTradeNoteUsecase(mockRepo)
	ctx := context.Background()

	req := &dto.TradeNoteRequest{
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "購入理由: 業績好調",
	}

	expectedError := errors.New("database error")
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.TradeNote")).Return(expectedError)

	result, err := usecase.CreateTradeNote(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestTradeNoteUsecase_GetTradeNote(t *testing.T) {
	mockRepo := new(MockTradeNoteRepository)
	usecase := NewTradeNoteUsecase(mockRepo)
	ctx := context.Background()

	expectedID := uint(1)
	expectedTradeNote := &entity.TradeNote{
		ID:        expectedID,
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "購入理由: 業績好調",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", ctx, expectedID).Return(expectedTradeNote, nil)

	result, err := usecase.GetTradeNote(ctx, expectedID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTradeNote.ID, result.ID)
	assert.Equal(t, expectedTradeNote.TradeDate, result.TradeDate)
	assert.Equal(t, expectedTradeNote.Symbol, result.Symbol)
	assert.Equal(t, expectedTradeNote.Content, result.Content)
	mockRepo.AssertExpectations(t)
}

func TestTradeNoteUsecase_GetAllTradeNotes(t *testing.T) {
	mockRepo := new(MockTradeNoteRepository)
	usecase := NewTradeNoteUsecase(mockRepo)
	ctx := context.Background()

	expectedTradeNotes := []*entity.TradeNote{
		{
			ID:        1,
			TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
			Symbol:    "AAPL",
			Content:   "購入",
		},
		{
			ID:        2,
			TradeDate: time.Date(2024, 6, 22, 9, 0, 0, 0, time.UTC),
			Symbol:    "GOOGL",
			Content:   "売却",
		},
	}

	mockRepo.On("GetAll", ctx).Return(expectedTradeNotes, nil)

	result, err := usecase.GetAllTradeNotes(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.TradeNotes, 2)
	assert.Equal(t, "AAPL", result.TradeNotes[0].Symbol)
	assert.Equal(t, "GOOGL", result.TradeNotes[1].Symbol)
	mockRepo.AssertExpectations(t)
}

func TestTradeNoteUsecase_UpdateTradeNote(t *testing.T) {
	mockRepo := new(MockTradeNoteRepository)
	usecase := NewTradeNoteUsecase(mockRepo)
	ctx := context.Background()

	targetID := uint(1)
	existingTradeNote := &entity.TradeNote{
		ID:        targetID,
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "購入理由: 業績好調",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &dto.TradeNoteRequest{
		TradeDate: time.Date(2024, 6, 24, 11, 0, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "更新された内容",
	}

	mockRepo.On("GetByID", ctx, targetID).Return(existingTradeNote, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.TradeNote")).Return(nil)

	result, err := usecase.UpdateTradeNote(ctx, targetID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, targetID, result.ID)
	assert.Equal(t, req.TradeDate, result.TradeDate)
	assert.Equal(t, req.Symbol, result.Symbol)
	assert.Equal(t, req.Content, result.Content)
	mockRepo.AssertExpectations(t)
}

func TestTradeNoteUsecase_DeleteTradeNote(t *testing.T) {
	mockRepo := new(MockTradeNoteRepository)
	usecase := NewTradeNoteUsecase(mockRepo)
	ctx := context.Background()

	targetID := uint(1)

	mockRepo.On("Delete", ctx, targetID).Return(nil)

	err := usecase.DeleteTradeNote(ctx, targetID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTradeNoteUsecase_DeleteTradeNote_Error(t *testing.T) {
	mockRepo := new(MockTradeNoteRepository)
	usecase := NewTradeNoteUsecase(mockRepo)
	ctx := context.Background()

	targetID := uint(1)
	expectedError := errors.New("database error")

	mockRepo.On("Delete", ctx, targetID).Return(expectedError)

	err := usecase.DeleteTradeNote(ctx, targetID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}