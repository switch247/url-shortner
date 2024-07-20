// models/url.go
package models

import "time"

type URL struct {
	ID         string `gorm:"primaryKey;size:10"`
	LongURL    string `gorm:"column:long_url;not null"`
	CreatedAt  time.Time
	ExpiresAt  *time.Time
	ClickCount int64 `gorm:"default:0"`
}

func (URL) TableName() string { return "urls" }
