// models/url.go
package models

import "time"

type Click struct {
	ID        uint   `gorm:"primaryKey"`
	ShortCode string `gorm:"index;size:10"`
	IP        string
	Country   string
	City      string
	UserAgent string
	Referrer  string
	Device    string
	OS        string
	Browser   string
	Utms      string
	CreatedAt time.Time
}

func (Click) TableName() string { return "clicks" }
