package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"ness/internal/application/dto"
	"ness/internal/application/usecase"
	"ness/internal/presentation/request"
)

type TradeNoteHandler struct {
	tradeNoteUsecase *usecase.TradeNoteUsecase
}

func NewTradeNoteHandler(tradeNoteUsecase *usecase.TradeNoteUsecase) *TradeNoteHandler {
	return &TradeNoteHandler{
		tradeNoteUsecase: tradeNoteUsecase,
	}
}

func (h *TradeNoteHandler) GetTradeNotes(c echo.Context) error {
	ctx := c.Request().Context()
	
	tradeNotes, err := h.tradeNoteUsecase.GetAllTradeNotes(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get trade notes",
		})
	}

	return c.JSON(http.StatusOK, tradeNotes)
}

func (h *TradeNoteHandler) CreateTradeNote(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req dto.TradeNoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := request.ValidateRequest(c, &req); err != nil {
		return err
	}

	tradeNote, err := h.tradeNoteUsecase.CreateTradeNote(ctx, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create trade note",
		})
	}

	return c.JSON(http.StatusCreated, tradeNote)
}

func (h *TradeNoteHandler) GetTradeNote(c echo.Context) error {
	ctx := c.Request().Context()
	
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID parameter",
		})
	}

	tradeNote, err := h.tradeNoteUsecase.GetTradeNote(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Trade note not found",
		})
	}

	return c.JSON(http.StatusOK, tradeNote)
}

func (h *TradeNoteHandler) UpdateTradeNote(c echo.Context) error {
	ctx := c.Request().Context()
	
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID parameter",
		})
	}

	var req dto.TradeNoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := request.ValidateRequest(c, &req); err != nil {
		return err
	}

	tradeNote, err := h.tradeNoteUsecase.UpdateTradeNote(ctx, uint(id), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update trade note",
		})
	}

	return c.JSON(http.StatusOK, tradeNote)
}

func (h *TradeNoteHandler) DeleteTradeNote(c echo.Context) error {
	ctx := c.Request().Context()
	
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID parameter",
		})
	}

	if err := h.tradeNoteUsecase.DeleteTradeNote(ctx, uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete trade note",
		})
	}

	return c.NoContent(http.StatusNoContent)
}