package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ness/internal/application/dto"
)

import (
	"ness/internal/application/usecase"
)

// MockTradeNoteUsecase is a mock implementation of TradeNoteUsecaseInterface
type MockTradeNoteUsecase struct {
	mock.Mock
}

func (m *MockTradeNoteUsecase) CreateTradeNote(ctx context.Context, req *dto.TradeNoteRequest) (*dto.TradeNoteResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TradeNoteResponse), args.Error(1)
}

func (m *MockTradeNoteUsecase) GetTradeNote(ctx context.Context, id uint) (*dto.TradeNoteResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TradeNoteResponse), args.Error(1)
}

func (m *MockTradeNoteUsecase) GetAllTradeNotes(ctx context.Context) (*dto.TradeNoteListResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TradeNoteListResponse), args.Error(1)
}

func (m *MockTradeNoteUsecase) UpdateTradeNote(ctx context.Context, id uint, req *dto.TradeNoteRequest) (*dto.TradeNoteResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TradeNoteResponse), args.Error(1)
}

func (m *MockTradeNoteUsecase) DeleteTradeNote(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Ensure MockTradeNoteUsecase implements the interface
var _ usecase.TradeNoteUsecaseInterface = (*MockTradeNoteUsecase)(nil)

func TestTradeNoteHandler_GetTradeNotes(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	expectedResponse := &dto.TradeNoteListResponse{
		TradeNotes: []dto.TradeNoteResponse{
			{
				ID:        1,
				TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
				Symbol:    "AAPL",
				Content:   "購入",
			},
		},
		Total: 1,
	}

	mockUsecase.On("GetAllTradeNotes", mock.Anything).Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetTradeNotes(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.TradeNoteListResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Total)
	assert.Len(t, response.TradeNotes, 1)
	mockUsecase.AssertExpectations(t)
}

func TestTradeNoteHandler_CreateTradeNote(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	requestBody := dto.TradeNoteRequest{
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "購入理由: 業績好調",
	}

	expectedResponse := &dto.TradeNoteResponse{
		ID:        1,
		TradeDate: requestBody.TradeDate,
		Symbol:    requestBody.Symbol,
		Content:   requestBody.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockUsecase.On("CreateTradeNote", mock.Anything, mock.AnythingOfType("*dto.TradeNoteRequest")).Return(expectedResponse, nil)

	bodyBytes, _ := json.Marshal(requestBody)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateTradeNote(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.TradeNoteResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, requestBody.Symbol, response.Symbol)
	mockUsecase.AssertExpectations(t)
}

func TestTradeNoteHandler_CreateTradeNote_InvalidJSON(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateTradeNote(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTradeNoteHandler_GetTradeNote(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	expectedID := uint(1)
	expectedResponse := &dto.TradeNoteResponse{
		ID:        expectedID,
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "購入理由: 業績好調",
	}

	mockUsecase.On("GetTradeNote", mock.Anything, expectedID).Return(expectedResponse, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/notes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetTradeNote(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.TradeNoteResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, response.ID)
	mockUsecase.AssertExpectations(t)
}

func TestTradeNoteHandler_GetTradeNote_InvalidID(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/notes/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := handler.GetTradeNote(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTradeNoteHandler_UpdateTradeNote(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	targetID := uint(1)
	requestBody := dto.TradeNoteRequest{
		TradeDate: time.Date(2024, 6, 24, 11, 0, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "更新された内容",
	}

	expectedResponse := &dto.TradeNoteResponse{
		ID:        targetID,
		TradeDate: requestBody.TradeDate,
		Symbol:    requestBody.Symbol,
		Content:   requestBody.Content,
		UpdatedAt: time.Now(),
	}

	mockUsecase.On("UpdateTradeNote", mock.Anything, targetID, mock.AnythingOfType("*dto.TradeNoteRequest")).Return(expectedResponse, nil)

	bodyBytes, _ := json.Marshal(requestBody)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/notes/1", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.UpdateTradeNote(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.TradeNoteResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, targetID, response.ID)
	assert.Equal(t, requestBody.Content, response.Content)
	mockUsecase.AssertExpectations(t)
}

func TestTradeNoteHandler_DeleteTradeNote(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	targetID := uint(1)

	mockUsecase.On("DeleteTradeNote", mock.Anything, targetID).Return(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeleteTradeNote(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUsecase.AssertExpectations(t)
}

func TestTradeNoteHandler_DeleteTradeNote_Error(t *testing.T) {
	mockUsecase := new(MockTradeNoteUsecase)
	handler := NewTradeNoteHandler(mockUsecase)

	targetID := uint(1)
	expectedError := errors.New("database error")

	mockUsecase.On("DeleteTradeNote", mock.Anything, targetID).Return(expectedError)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeleteTradeNote(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUsecase.AssertExpectations(t)
}