package http

import (
	"url-checker/internal/service"
	"url-checker/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(checkerService *service.CheckerService) *gin.Engine {
	r := gin.Default()

	pdfGenerator := service.NewGenerator()
	handler := handlers.NewHandler(checkerService, pdfGenerator)

	api := r.Group("/api/v1")
	{
		api.POST("/check", handler.CheckLinks)
		api.POST("/report", handler.GenerateReport)
	}

	return r
}
