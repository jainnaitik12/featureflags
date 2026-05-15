package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"featureflags/audit-service/internal/controller"
	"featureflags/audit-service/internal/repository"
	"featureflags/audit-service/internal/service"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer pool.Close()

	repo := repository.NewAuditPostgres(pool)
	if err := repo.EnsureSchema(ctx); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	svc := service.NewAuditService(repo)
	ctrl := controller.NewAuditController(svc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	ctrl.RegisterRoutes(r)

	log.Println("audit-service listening on :8084")
	if err := r.Run(":8084"); err != nil {
		log.Fatal(err)
	}
}
