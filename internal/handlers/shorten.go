// internal/handlers/shorten.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"tinyurl-analytics/internal/config"
	"tinyurl-analytics/internal/db"
	"tinyurl-analytics/internal/models"
	"tinyurl-analytics/internal/service"

	"github.com/go-chi/render"
)

type ShortenRequest struct {
	URL string `json:"url" form:"url"`
}

type ShortenResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
	AppName   string `json:"-"`
	AppDomain string `json:"-"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	// Accept both JSON and form data
	contentType := r.Header.Get("Content-Type")

	if r.Method == "POST" && (contentType == "" || contains(contentType, "application/x-www-form-urlencoded")) {
		// HTMX form → parse form
		if err := r.ParseForm(); err != nil {
			renderError(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		req.URL = r.FormValue("url")
	} else {
		// JSON API request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			renderError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	}

	if req.URL == "" {
		renderError(w, "Missing url", http.StatusBadRequest)
		return
	}

	code := service.GenerateShortCode()

	urlModel := models.URL{
		ID:      code,
		LongURL: req.URL,
	}

	if err := db.GORM.Create(&urlModel).Error; err != nil {
		log.Println("DB error:", err)
		renderError(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	db.Redis.Set(db.Ctx, "url:"+code, req.URL, 0)

	resp := ShortenResponse{
		ShortURL:  fmt.Sprintf("%s/%s", config.AppDomain, code),
		ShortCode: code,
		LongURL:   req.URL,
		AppName:   config.AppName,
		AppDomain: config.AppDomain,
	}

	// Check if it's an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		if ResultTmpl == nil {
			if err := InitTemplates(); err != nil {
				log.Println("InitTemplates error:", err)
				renderError(w, "Template init error", http.StatusInternalServerError)
				return
			}
		}

		// Return HTML for HTMX
		w.Header().Set("Content-Type", "text/html")
		if err := ResultTmpl.ExecuteTemplate(w, "result", resp); err != nil {
			log.Println("Template execution error:", err)
			renderError(w, "Failed to render", http.StatusInternalServerError)
			return
		}
		return
	}

	render.JSON(w, r, resp)
}

func renderError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s[:len(substr)+1] == substr+"/" || s[:len(substr)+1] == substr+";")
}
