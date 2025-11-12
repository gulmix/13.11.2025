package http

import (
	"url-checker/internal/service"
	"url-checker/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(svc *service.Service) *gin.Engine {
	r := gin.Default()
	h := handlers.NewHandler(svc)

	api := r.Group("/api/v1")
	{
		api.POST("/check", h.CheckLinks)
		api.POST("/report", h.GetReport)
	}

	return r
}
