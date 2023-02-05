package handler

import (
	"aquafarm-management/app/model"
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PondHandler menangani permintaan HTTP untuk item
type PondHandler struct {
	pondUsecase *usecase.PondUsecase
}

// NewPondHandler membuat instance baru ItemHandler
func NewPondHandler(iu *usecase.PondUsecase) *PondHandler {
	return &PondHandler{
		pondUsecase: iu,
	}
}

// Fetch
// @Summary Get All Ponds Data
// @Description Get All Ponds Data by Joining Farm Table.
// @Tags Ponds
// @Accept  json
// @Produce  json
// @Success 200 {object} response.HTTPResponseDataPonds
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/pond/ [get]
// @Security APIKeyHeader
// @Header 200 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *PondHandler) Fetch(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	farms, err := h.pondUsecase.Fetch()
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
// @Summary Get One Pond Data
// @Description Get one Pond Data by Joining Farm Table.
// @Tags Ponds
// @Accept  json
// @Produce  json
// @Success 200 {object} response.HTTPResponseDataPond
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Failure 404 {object} response.HTTPResponseAction
// @Param farmId path string true "Farm ID"
// @Param pondId path string true "Pond ID"
// @Router /v1/pond/{farmId}/{pondId} [get]
// @Security APIKeyHeader
// @Header 200 {string} Authorization "apiKey"
// @Header 404 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *PondHandler) GetById(c *gin.Context) {
	// mengambil ID item dari URL
	farmId := c.Param("farm_id")
	pondId := c.Param("pond_id")
	// memanggil usecase untuk mengambil item
	item, isNotFound, err := h.pondUsecase.GetById(farmId, pondId)
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
			"message": "Data Pond '" + pondId + "' of Farm '" + farmId + "' Not Found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   item,
	})
}

// Store
// @Summary Insert new Pond data
// @Description Insert one Pond Data by send the FarmdId.
// @Tags Ponds
// @Accept  json
// @Produce  json
// @Param pond body model.PondStore true "Required Data to Insert Pond"
// @Success 201 {object} response.HTTPResponseAction
// @Failure 409 {object} response.HTTPResponseAction
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/pond [post]
// @Security APIKeyHeader
// @Header 201 {string} Authorization "apiKey"
// @Header 409 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *PondHandler) Store(c *gin.Context) {
	// mengambil data item dari permintaan HTTP
	var pond model.Pond
	err := c.BindJSON(&pond)
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if pond.FarmID == "" || pond.ID == "" || pond.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Please read documentation, system need more data",
		})
		return
	}

	var isDuplicateEntry bool
	// memanggil usecase untuk membuat item baru
	_, isDuplicateEntry, err = h.pondUsecase.Store(&pond)
	if isDuplicateEntry {
		c.JSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Data Pond '" + pond.ID + "' of Farm '" + pond.FarmID + "' still exist",
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
// @Summary Update exist Pond data or Insert new
// @Description Update exist Pond data or Insert new if not exist.
// @Tags Ponds
// @Accept  json
// @Produce  json
// @Param pond body model.PondStore true "Required Data to Update or Insert Pond"
// @Success 201 {object} response.HTTPResponseAction "Create Data"
// @Success 202 {object} response.HTTPResponseAction "Update Data"
// @Failure 500 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Router /v1/pond [put]
// @Security APIKeyHeader
// @Header 201 {string} Authorization "apiKey"
// @Header 202 {string} Authorization "apiKey"
// @Header 500 {string} Authorization "apiKey"
func (h *PondHandler) UpdateById(c *gin.Context) {
	// mengambil body JSON farm dari Body Post
	var pond model.Pond
	err := c.BindJSON(&pond)
	if err != nil {
		// menangani error
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// memanggil usecase untuk memperbarui item
	isCreateData, err := h.pondUsecase.UpdateById(&pond)
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
			"message": "Create data Pond '" + pond.ID + "' of Farm '" + pond.FarmID + "' successfully",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  http.StatusAccepted,
		"message": "Update data Pond '" + pond.ID + "' of Farm '" + pond.FarmID + "' successfully",
	})
}

// SoftDeleteById
// @Summary (soft) Delete exist Pond data
// @Description Change Flagger (is_deleted) value on Database to TRUE.
// @Tags Ponds
// @Accept  json
// @Produce  json
// @Success 202 {object} response.HTTPResponseDataPond
// @Failure 502 {object} response.HTTPResponseAction
// @Failure 401 {object} response.HTTPResponseAction
// @Failure 404 {object} response.HTTPResponseAction
// @Param farmId path string true "Farm ID"
// @Param pondId path string true "Pond ID"
// @Router /v1/pond/{farmId}/{pondId} [delete]
// @Security APIKeyHeader
// @Header 202 {string} Authorization "apiKey"
// @Header 404 {string} Authorization "apiKey"
// @Header 502 {string} Authorization "apiKey"
func (h *PondHandler) SoftDeleteById(c *gin.Context) {
	farmId := c.Param("farm_id")
	pondId := c.Param("pond_id")
	isDataNotFound, err := h.pondUsecase.SoftDeleteById(farmId, pondId)
	if isDataNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Pond '" + pondId + "' of Farm '" + farmId + "' already not found",
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
		"message": "Delete Pond '" + pondId + "' of Farm '" + farmId + "' successfully",
	})
}
