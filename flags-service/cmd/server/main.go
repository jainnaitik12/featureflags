package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"featureflags/flags-service/internal/controller"
	"featureflags/flags-service/internal/repository"
	"featureflags/flags-service/internal/service"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is required")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("invalid REDIS_URL: %v", err)
	}

	rdb := redis.NewClient(opts)
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	repo := repository.NewFlagsRedis(rdb)
	svc := service.NewFlagService(repo)
	flagsCtrl := controller.NewFlagsController(svc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	flagsCtrl.RegisterRoutes(r)

	log.Println("flags-service listening on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
