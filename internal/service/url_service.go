// internal/service/url_service.go
package service

import (
	"crypto/rand"

	"tinyurl-analytics/internal/db"
	"tinyurl-analytics/internal/models"
)

var counter uint64 = 100000 // Start from a high number so short codes look nice

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const length = 7 // 62^7 ≈ 3.5 trillion possible codes

// GenerateShortCode - Collision-resistant, fast, sequential-ish
func GenerateShortCode() string {
	for {
		// Option 1: Fully random (simpler, works great)
		b := make([]byte, length)
		rand.Read(b)
		for i := range b {
			b[i] = charset[b[i]%62]
		}
		code := string(b)

		// Check collision (extremely rare)
		var exists bool
		db.GORM.Model(&models.URL{}).
			Select("count(*) > 0").
			Where("id = ?", code).
			Scan(&exists)

		if !exists {
			return code
		}
	}
}

// BONUS: Want prettier sequential codes like bit.ly? Use this instead:
// func GenerateShortCode() string {
//     id := atomic.AddUint64(&counter, 1)
//     return base62.Encode(id)
// }
