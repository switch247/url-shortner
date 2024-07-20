package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"tinyurl-analytics/internal/config"
	"tinyurl-analytics/internal/db"
	"tinyurl-analytics/internal/models"

	"github.com/go-chi/chi/v5"
)

// CONFIG
// const AppDomain = "http://localhost:8080" // Moved to config package

// TEMPLATE FUNCTIONS
func defaultFunc(arg, def string) string {
	if arg == "" {
		return def
	}
	return arg
}

func sliceFunc(s string, start, end int) string {
	if start > len(s) {
		start = len(s)
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		start = end
	}
	return s[start:end]
}

// DATA STRUCTURES
type StatsData struct {
	AppDomain  string
	AppName    string
	ShortCode  string
	LongURL    string
	ClickCount int64
	Clicks     []models.Click
}

// GLOBAL TEMPLATES
var (
	IndexTmpl        *template.Template
	DashboardTmpl    *template.Template
	AllStatsTmpl     *template.Template
	StatsPartialTmpl *template.Template
	ResultTmpl       *template.Template
)

// ✅ Call this once on server startup
func InitTemplates() error {
	var err error
	funcMap := template.FuncMap{
		"default": defaultFunc,
		"slice":   sliceFunc,
		"printf":  fmt.Sprintf,
	}

	// Helper to parse a set of files
	parse := func(files ...string) (*template.Template, error) {
		return template.New("").Funcs(funcMap).ParseFiles(files...)
	}

	// 1. Home Page
	IndexTmpl, err = parse("web/templates/base.html", "web/templates/index.html")
	if err != nil {
		return fmt.Errorf("failed to parse index templates: %w", err)
	}

	// 2. Dashboard (Single Stats)
	DashboardTmpl, err = parse("web/templates/base.html", "web/templates/dashboard.html")
	if err != nil {
		return fmt.Errorf("failed to parse dashboard templates: %w", err)
	}

	// 3. All Stats Page
	AllStatsTmpl, err = parse("web/templates/base.html", "web/templates/all_stats.html")
	if err != nil {
		return fmt.Errorf("failed to parse all_stats templates: %w", err)
	}

	// 4. Stats Partial (HTMX) - No base layout
	StatsPartialTmpl, err = parse("web/templates/stats_partial.html")
	if err != nil {
		return fmt.Errorf("failed to parse stats_partial template: %w", err)
	}

	// 5. Result Partial (HTMX) - No base layout
	ResultTmpl, err = parse("web/templates/result.html")
	if err != nil {
		return fmt.Errorf("failed to parse result template: %w", err)
	}

	return nil
}

// --- HANDLERS ---

// GET /{code}+ → Stats page for single URL
func StatsPage(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if db.GORM == nil {
		http.Error(w, "database not initialized", http.StatusInternalServerError)
		return
	}

	if DashboardTmpl == nil {
		if err := InitTemplates(); err != nil {
			log.Printf("InitTemplates error: %v", err)
			http.Error(w, "template init error", http.StatusInternalServerError)
			return
		}
	}

	var url models.URL
	if err := db.GORM.Where("id = ?", code).First(&url).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	var clickCount int64
	if err := db.GORM.Model(&models.Click{}).Where("short_code = ?", code).Count(&clickCount).Error; err != nil {
		http.Error(w, "failed to fetch click count", http.StatusInternalServerError)
		return
	}

	data := StatsData{
		AppDomain:  config.AppDomain,
		AppName:    config.AppName,
		ShortCode:  code,
		LongURL:    url.LongURL,
		ClickCount: clickCount,
	}

	if err := DashboardTmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /stats/{code} → HTMX partial
func StatsPartial(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var count int64
	if db.GORM == nil {
		http.Error(w, "database not initialized", http.StatusInternalServerError)
		return
	}

	if StatsPartialTmpl == nil {
		if err := InitTemplates(); err != nil {
			log.Printf("InitTemplates error: %v", err)
			http.Error(w, "template init error", http.StatusInternalServerError)
			return
		}
	}

	if res := db.GORM.Model(&models.Click{}).
		Where("short_code = ?", code).
		Count(&count); res.Error != nil {
		http.Error(w, "failed to fetch click count", http.StatusInternalServerError)
		return
	}

	var clicks []models.Click
	if res := db.GORM.Where("short_code = ?", code).
		Order("created_at DESC").
		Limit(50).
		Find(&clicks); res.Error != nil {
		http.Error(w, "failed to fetch clicks", http.StatusInternalServerError)
		return
	}

	data := StatsData{
		AppDomain:  config.AppDomain,
		AppName:    config.AppName,
		ShortCode:  code,
		ClickCount: count,
		Clicks:     clicks,
	}

	if err := StatsPartialTmpl.ExecuteTemplate(w, "stats_partial.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /all_stats → Show stats for all URLs
func AllStatsPage(w http.ResponseWriter, r *http.Request) {
	if AllStatsTmpl == nil {
		if err := InitTemplates(); err != nil {
			log.Printf("InitTemplates error: %v", err)
			http.Error(w, "template init error", http.StatusInternalServerError)
			return
		}
	}

	var urls []models.URL
	if err := db.GORM.Order("created_at DESC").Find(&urls).Error; err != nil {
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched %d URLs for all stats", len(urls))

	var allStats []URLWithClicks
	for i, u := range urls {
		var clicks []models.Click
		if res := db.GORM.Where("short_code = ?", u.ID).
			Order("created_at DESC").
			Limit(50).
			Find(&clicks); res.Error != nil {
			log.Printf("failed to fetch clicks for %s: %v", u.ID, res.Error)
			clicks = nil
		}

		var clickCount int64
		if err := db.GORM.Model(&models.Click{}).Where("short_code = ?", u.ID).Count(&clickCount).Error; err != nil {
			log.Printf("failed to count clicks for %s: %v", u.ID, err)
			clickCount = 0
		}
		urls[i].ClickCount = clickCount

		allStats = append(allStats, URLWithClicks{
			URL:    urls[i],
			Clicks: clicks,
		})
	}
	log.Printf("Prepared stats for %d URLs", len(allStats))

	err := AllStatsTmpl.ExecuteTemplate(w, "all_stats.html", struct {
		AppDomain string
		AppName   string
		Stats     []URLWithClicks
	}{
		AppDomain: config.AppDomain,
		AppName:   config.AppName,
		Stats:     allStats,
	})

	if err != nil {
		log.Printf("ExecuteTemplate error: %v", err)
		os.WriteFile("c:\\BackUp\\web-projects\\Go Lang\\url-shortner\\error.log", []byte(err.Error()), 0644)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type URLWithClicks struct {
	models.URL
	Clicks []models.Click
}

// GET / → Home page (rendered template)
func IndexPage(w http.ResponseWriter, r *http.Request) {
	if IndexTmpl == nil {
		if err := InitTemplates(); err != nil {
			log.Printf("InitTemplates error: %v", err)
			http.Error(w, "template init error", http.StatusInternalServerError)
			return
		}
	}

	// pass minimal app info used by templates
	data := struct {
		AppDomain string
		AppName   string
	}{
		AppDomain: config.AppDomain,
		AppName:   config.AppName,
	}

	if err := IndexTmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
