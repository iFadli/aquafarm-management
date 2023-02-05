package handler

import (
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
)

// StatHandler menangani permintaan HTTP untuk item
type StatHandler struct {
	statUsecase *usecase.StatUsecase
}

// NewStatHandler membuat instance baru StatHandler
func NewStatHandler(iu *usecase.StatUsecase) *StatHandler {
	return &StatHandler{
		statUsecase: iu,
	}
}

// Fetch mengambil semua item
func (h *StatHandler) Fetch(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	stats, err := h.statUsecase.Fetch()
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"data":   stats.StatisticsData,
	})
}

// FetchLogs mengambil semua item
func (h *StatHandler) FetchLogs(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	logs, err := h.statUsecase.FetchLogs()
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"data":   logs,
	})
}
