package router

import (
	"net/http"
	"os"
	"time"

	"github.com/gadzz/scopepay/internal/config"
	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gadzz/scopepay/internal/handlers"
	"github.com/gadzz/scopepay/internal/service"
	gormstore "github.com/gadzz/scopepay/internal/store/gorm"
)

func New(cfg config.Config, database *gorm.DB) *gin.Engine {
	r := gin.New()

	// Middleware setup
	r.Use(httpx.RequestID(), gin.Logger(), httpx.Recovery())

	// Databse for project setup
	projectStore := gormstore.NewProjectStore(database)
	projectSvc := service.NewProjectService(projectStore)
	projectHandler := handlers.NewProjectHandler(projectSvc)

	origins := []string{"http://localhost:5173"}
	if v := os.Getenv("CORS_ORIGINS"); v != "" {
		origins = []string{}
		for _, o := range splitAndTrim(v) {
			origins = append(origins, o)
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     cfg.DBDriver,
		})
	})

	r.POST("/projects", projectHandler.Create)

	return r
}

func splitAndTrim(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := trimSpace(s[start:i])
			if part != "" {
				out = append(out, part)
			}
			start = i + 1
		}
	}
	return out
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
