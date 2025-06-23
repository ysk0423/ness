package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ness/internal/domain/entity"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return db, mock
}

func TestTradeNoteRepositoryImpl_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTradeNoteRepository(db)
	ctx := context.Background()

	tradeNote := &entity.TradeNote{
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "購入理由: 業績好調",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "trade_notes"`).
		WithArgs(
			tradeNote.TradeDate,
			tradeNote.Symbol,
			tradeNote.Content,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(ctx, tradeNote)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeNoteRepositoryImpl_GetByID(t *testing.T) {
	db, mock := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTradeNoteRepository(db)
	ctx := context.Background()

	expectedID := uint(1)
	expectedTradeDate := time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC)
	expectedSymbol := "AAPL"
	expectedContent := "購入理由: 業績好調"

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "trade_date", "symbol", "content"}).
		AddRow(expectedID, time.Now(), time.Now(), nil, expectedTradeDate, expectedSymbol, expectedContent)

	mock.ExpectQuery(`SELECT \* FROM "trade_notes"`).
		WithArgs(expectedID, 1).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, expectedID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedID, result.ID)
	assert.Equal(t, expectedTradeDate, result.TradeDate)
	assert.Equal(t, expectedSymbol, result.Symbol)
	assert.Equal(t, expectedContent, result.Content)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeNoteRepositoryImpl_GetByID_NotFound(t *testing.T) {
	db, mock := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTradeNoteRepository(db)
	ctx := context.Background()

	expectedID := uint(999)

	mock.ExpectQuery(`SELECT \* FROM "trade_notes"`).
		WithArgs(expectedID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.GetByID(ctx, expectedID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeNoteRepositoryImpl_GetAll(t *testing.T) {
	db, mock := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTradeNoteRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "trade_date", "symbol", "content"}).
		AddRow(1, time.Now(), time.Now(), nil, time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC), "AAPL", "購入").
		AddRow(2, time.Now(), time.Now(), nil, time.Date(2024, 6, 22, 9, 0, 0, 0, time.UTC), "GOOGL", "売却")

	mock.ExpectQuery(`SELECT \* FROM "trade_notes"`).
		WillReturnRows(rows)

	result, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "AAPL", result[0].Symbol)
	assert.Equal(t, "GOOGL", result[1].Symbol)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeNoteRepositoryImpl_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTradeNoteRepository(db)
	ctx := context.Background()

	tradeNote := &entity.TradeNote{
		ID:        1,
		TradeDate: time.Date(2024, 6, 23, 10, 30, 0, 0, time.UTC),
		Symbol:    "AAPL",
		Content:   "更新された内容",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "trade_notes"`).
		WithArgs(
			tradeNote.TradeDate,
			tradeNote.Symbol,
			tradeNote.Content,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			tradeNote.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, tradeNote)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTradeNoteRepositoryImpl_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTradeNoteRepository(db)
	ctx := context.Background()

	targetID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "trade_notes" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), targetID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, targetID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}