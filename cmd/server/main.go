package main

import (
	"log"
	"net/http"
	"os"
	"tinyurl-analytics/internal/analytics"
	"tinyurl-analytics/internal/config"
	"tinyurl-analytics/internal/db"
	"tinyurl-analytics/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	db.Init()
	tracker := analytics.NewTracker()
	// ✅ Initialize HTML templates
	handlers.InitTemplates()
	// set logger to log to logs file
	logFile, err := os.OpenFile("logs/server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handlers.IndexPage)
	r.Post("/shorten", handlers.ShortenURL)
	r.Get("/{code}", handlers.Redirect(tracker))
	r.Get("/stats/{code}", handlers.StatsPartial) // helper partial for HTMX
	r.Get("/all_stats", handlers.AllStatsPage)    // all stats page
	r.Get("/{code}+", handlers.StatsPage)         // stats page

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))).ServeHTTP(w, r)
	})

	log.Println("Server running on :" + config.Port)
	http.ListenAndServe(":"+config.Port, r)
}
