package handlers

import (
	"net/http"
	"time"
	"tinyurl-analytics/internal/analytics"
	"tinyurl-analytics/internal/db"
	"tinyurl-analytics/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func Redirect(tracker *analytics.Tracker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`<html><body><h1>Not Found</h1><p>The requested short URL does not exist.</p></body></html>`))
			return
		}

		longURL, err := db.Redis.Get(db.Ctx, "url:"+code).Result()
		if err == redis.Nil {
			var url models.URL
			if err := db.GORM.Select("long_url").Where("id = ?", code).First(&url).Error; err != nil {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`<html><body><h1>Not Found</h1><p>The requested short URL does not exist.</p></body></html>`))
				return
			}
			longURL = url.LongURL
			db.Redis.Set(db.Ctx, "url:"+code, longURL, 24*time.Hour)
		}

		// Track click
		go tracker.Track(r, code)

		http.Redirect(w, r, longURL, http.StatusFound)
	}
}
