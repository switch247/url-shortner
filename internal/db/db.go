// internal/db/db.go
package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GORM *gorm.DB
var Redis *redis.Client
var Ctx = context.Background() // ← EXPORTED NOW (uppercase = public)

func Init() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	var err error
	GORM, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Postgres error:", err)
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Redis URL parse error:", err)
	}

	Redis = redis.NewClient(opt)

	if err := Redis.Ping(Ctx).Err(); err != nil {
		log.Fatal("Redis connection error:", err)
	}

	log.Println("Database & Redis connected")
}
