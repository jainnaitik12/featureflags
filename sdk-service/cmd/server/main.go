package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"featureflags/sdk-service/internal/controller"
	"featureflags/sdk-service/internal/metrics"
	"featureflags/sdk-service/internal/repository"
	"featureflags/sdk-service/internal/service"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is required")
	}

	metrics.Register()

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("invalid REDIS_URL: %v", err)
	}
	rdb := redis.NewClient(opts)
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	repo := repository.NewFlagsRedisRead(rdb)
	svc := service.NewEvalService(repo)
	ctrl := controller.NewSDKController(svc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	ctrl.RegisterRoutes(r)

	log.Println("sdk-service listening on :8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}
