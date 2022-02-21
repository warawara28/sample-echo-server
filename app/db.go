package app

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database データベースの接続を管理
type Database struct {
	db *gorm.DB
}

type dbLogger struct {
	logger Logger
}

func (l *dbLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	return l
}
func (l *dbLogger) Info(_ context.Context, fmt string, args ...interface{}) {
	l.logger.Info(fmt, args...)
}
func (l *dbLogger) Warn(_ context.Context, fmt string, args ...interface{}) {
	l.logger.Warning(fmt, args...)
}
func (l *dbLogger) Error(_ context.Context, fmt string, args ...interface{}) {
	l.logger.Error(fmt, args...)
}
func (l *dbLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rowsAffected := fc()
	l.logger.Debug("begin:%v sql:%s affected:%d err:%v", begin, sql, rowsAffected, err)
}

// NewDatabase データベースの接続を管理する構造体を作成
func NewDatabase(dsn string, logger Logger) (*Database, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed gorm open error:%w", err)
	}
	gormDB.Logger = &dbLogger{logger: logger}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed gorm getDB error:%w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed gorm ping error:%w", err)
	}

	return &Database{
		db: gormDB,
	}, nil
}

// AutoMigrate 構造体からマイグレーションを行う
func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}

// ListBooks 書籍を全件取得
func (d *Database) ListBooks() ([]*Book, error) {
	books := []*Book{}
	if err := d.db.Find(&books).Error; err != nil {
		return nil, fmt.Errorf("failed ListBooks: %w", err)
	}
	return books, nil
}
