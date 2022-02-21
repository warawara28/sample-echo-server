package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Service APIのハンドラー実装
type Service struct {
	db *Database
}

// ErrorResponse エラー時のレスポンス用構造体
type ErrorResponse struct {
	Error string
}

// NewService APIのハンドラー実装を生成
func NewService(db *Database) *Service {
	return &Service{
		db: db,
	}
}

// CheckHealth ヘルスチェック用
func (srv *Service) CheckHealth() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

// ListBooks 全書籍データをJSON形式で返却
func (srv *Service) ListBooks() echo.HandlerFunc {
	return func(c echo.Context) error {
		books, err := srv.db.ListBooks()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{err.Error()})
		}
		return c.JSON(http.StatusOK, books)
	}
}
