package handler

import (
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
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

// Fetch
// @Summary Get Statistics Data
// @Description Get Statistics Data From Logs Middleware.
// @Tags Logs
// @Accept  json
// @Produce  json
// @Success 200 {object} response.HTTPResponseDataStatistics
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/logs/statistics [get]
// @Security APIKeyHeader
// @Header 200 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *StatHandler) Fetch(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	stats, err := h.statUsecase.Fetch()
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   stats.StatisticsData,
	})
}

// FetchLogs
// @Summary Get All Logs Data
// @Description Get All Logs Data From Logs Middleware.
// @Tags Logs
// @Accept  json
// @Produce  json
// @Success 200 {object} response.HTTPResponseDataLogs
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/logs [get]
// @Security APIKeyHeader
// @Header 200 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *StatHandler) FetchLogs(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	logs, err := h.statUsecase.FetchLogs()
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   logs,
	})
}
