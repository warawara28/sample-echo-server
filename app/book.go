package app

import "time"

// Book 書籍データ
type Book struct {
	ID          int64     `gorm:"type:int;primaryKey;autoIncrement;comment:書籍名"`
	Title       string    `gorm:"type:varchar(255);not null;comment:書籍名"`
	TotalPages  int64     `gorm:"type:int;not null;comment:総ページ数"`
	Rating      float64   `gorm:"type:decimal(4,2);not null;comment:出版日時"`
	OnSale      bool      `gorm:"type:tinyint(1);not null;default:1;comment:販売中=1,販売停止=0"`
	PublishedAt time.Time `gorm:"type:datetime;not null;comment:出版日時"`
	UpdatedAt   time.Time `gorm:"type:datetime;not null;autoUpdateTime;comment:レコード更新日時"`
	CreatedAt   time.Time `gorm:"type:datetime;not null;autoCreateTime;comment:レコード登録日時"`
}
