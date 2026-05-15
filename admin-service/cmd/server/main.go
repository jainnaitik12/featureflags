package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"featureflags/admin-service/internal/controller"
	"featureflags/admin-service/internal/repository"
	"featureflags/admin-service/internal/service"
)

func main() {
	flagsServiceURL := getEnv("FLAGS_SERVICE_URL", "http://flags-service:8081")
	auditServiceURL := getEnv("AUDIT_SERVICE_URL", "http://audit-service:8084")

	flagsUpstream := repository.NewFlagsUpstream(flagsServiceURL)
	auditUpstream := repository.NewAuditUpstream(auditServiceURL)
	svc := service.NewAdminService(flagsUpstream, auditUpstream)
	ctrl := controller.NewAdminController(svc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	ctrl.RegisterRoutes(r)

	log.Println("admin-service listening on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

func getEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
