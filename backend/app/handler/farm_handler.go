package handler

import (
	"aquafarm-management/app/model"
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FarmHandler menangani permintaan HTTP untuk item
type FarmHandler struct {
	farmUsecase *usecase.FarmUsecase
}

// NewFarmHandler membuat instance baru ItemHandler
func NewFarmHandler(iu *usecase.FarmUsecase) *FarmHandler {
	return &FarmHandler{
		farmUsecase: iu,
	}
}

// Fetch
// @Summary Get All Farms Data
// @Description Get All Farms Data but only is_deleted = false.
// @Tags Farms
// @Accept  json
// @Produce  json
// @Success 200 {object} response.HTTPResponseDataFarms
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/farm/ [get]
// @Security APIKeyHeader
// @Header 200 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *FarmHandler) Fetch(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	farms, err := h.farmUsecase.Fetch()
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
		"data":   farms,
	})
}

// GetById
// @Summary Get One Farm Data
// @Description Get one Farm Data but only is_deleted = false.
// @Tags Farms
// @Accept  json
// @Produce  json
// @Success 200 {object} response.HTTPResponseDataFarm
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Failure 404 {object} response.HTTPResponseAction
// @Param farmId path string true "Farm ID"
// @Router /v1/farm/{farmId} [get]
// @Security APIKeyHeader
// @Header 200 {string} Authorization "apiKey"
// @Header 404 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *FarmHandler) GetById(c *gin.Context) {
	// mengambil ID item dari URL
	farmId := c.Param("farm_id")
	// memanggil usecase untuk mengambil item
	item, isNotFound, err := h.farmUsecase.GetById(farmId)
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if isNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Farm '" + farmId + "' Not Found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   item,
	})
}

// Store
// @Summary Insert new Farm data
// @Description Insert one Farm Data.
// @Tags Farms
// @Accept  json
// @Produce  json
// @Param pond body model.FarmStore true "Required Data to Insert Farm"
// @Success 201 {object} response.HTTPResponseAction
// @Failure 409 {object} response.HTTPResponseAction
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/farm [post]
// @Security APIKeyHeader
// @Header 201 {string} Authorization "apiKey"
// @Header 409 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *FarmHandler) Store(c *gin.Context) {
	// mengambil data item dari permintaan HTTP
	var farm model.Farm
	err := c.BindJSON(&farm)
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if farm.ID == "" || farm.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Please read documentation, system need more data",
		})
		return
	}

	var isDuplicateEntry bool
	// memanggil usecase untuk membuat item baru
	_, isDuplicateEntry, err = h.farmUsecase.Store(&farm)
	if isDuplicateEntry {
		c.JSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Data Farm '" + farm.ID + "' still exist",
		})
		return
	}
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Create Data Successfully",
	})
}

// UpdateById
// @Summary Update exist Farm data or Insert new
// @Description Update exist Farm data or Insert new if not exist.
// @Tags Farms
// @Accept  json
// @Produce  json
// @Param pond body model.FarmStore true "Required Data to Update or Insert Farm"
// @Success 201 {object} response.HTTPResponseAction "Create Data"
// @Success 202 {object} response.HTTPResponseAction "Update Data"
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/farm [put]
// @Security APIKeyHeader
// @Header 201 {string} Authorization "apiKey"
// @Header 202 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *FarmHandler) UpdateById(c *gin.Context) {
	// mengambil body JSON farm dari Body Post
	var farm model.Farm
	err := c.BindJSON(&farm)
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// memanggil usecase untuk memperbarui item
	isCreateData, err := h.farmUsecase.UpdateById(&farm)
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}
	if isCreateData {
		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Create data Farm '" + farm.ID + "' successfully",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  http.StatusAccepted,
		"message": "Update data Farm '" + farm.ID + "' successfully",
	})
}

// SoftDeleteById
// @Summary (soft) Delete exist Farm data
// @Description Change Flagger (is_deleted) value on Database to TRUE.
// @Tags Farms
// @Accept  json
// @Produce  json
// @Success 202 {object} response.HTTPResponseDataPond
// @Failure 502 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Failure 404 {object} response.HTTPResponseAction
// @Param farmId path string true "Farm ID"
// @Param pondId path string true "Pond ID"
// @Router /v1/farm/{farmId} [delete]
// @Security APIKeyHeader
// @Header 202 {string} Authorization "apiKey"
// @Header 404 {string} Authorization "apiKey"
// @Header 502 {string} Authorization "apiKey"
func (h *FarmHandler) SoftDeleteById(c *gin.Context) {
	farmId := c.Param("farm_id")
	isDataNotFound, err := h.farmUsecase.SoftDeleteById(farmId)
	if isDataNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Farm '" + farmId + "' already not found",
		})
		return
	}
	if err != nil {
		// menangani error
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  http.StatusAccepted,
		"message": "Delete Farm '" + farmId + "' successfully",
	})
}
