package handlers

import (
	"net/http"
	"time"
	"url-checker/internal/models"
	"url-checker/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	checkerService *service.CheckerService
	pdfGenerator   *service.Generator
}

func NewHandler(checkerService *service.CheckerService, pdfGenerator *service.Generator) *Handler {
	return &Handler{
		checkerService: checkerService,
		pdfGenerator:   pdfGenerator,
	}
}

func (h *Handler) CheckLinks(c *gin.Context) {
	var req models.CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.checkerService.CheckLinks(req.Links)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < 50; i++ {
		if task, exists := h.checkerService.GetTask(task.ID); exists && task.Status == "completed" {
			c.JSON(http.StatusOK, models.CheckResponse{
				Links:   task.Results,
				LinksID: task.ID,
			})
			return
		}
		c.Request.Context().Done()
		select {
		case <-c.Request.Context().Done():
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
			return
		default:
			timer := time.NewTimer(100 * time.Millisecond)
			<-timer.C
		}
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "task accepted, still processing",
		"links_num": task.ID,
	})
}

func (h *Handler) GenerateReport(c *gin.Context) {
	var req models.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks := h.checkerService.GetTasks(req.LinksNum)
	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tasks found"})
		return
	}

	pdfData, err := h.pdfGenerator.GetReport(tasks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=report.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfData)
}
