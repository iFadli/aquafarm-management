package handler

import (
	"aquafarm-management/app/model"
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
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

// Fetch mengambil semua item
func (h *PondHandler) Fetch(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	farms, err := h.pondUsecase.Fetch()
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
		"data":   farms,
	})
}

// GetById mengambil Pond berdasarkan ID
func (h *PondHandler) GetById(c *gin.Context) {
	// mengambil ID item dari URL
	farmId := c.Param("farm_id")
	pondId := c.Param("pond_id")
	// memanggil usecase untuk mengambil item
	item, isNotFound, err := h.pondUsecase.GetById(farmId, pondId)
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	if isNotFound {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "pond data '" + farmId + "' tidak ada",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"data":   item,
	})
}

// Store membuat Pond baru
func (h *PondHandler) Store(c *gin.Context) {
	// mengambil data item dari permintaan HTTP
	var pond model.Pond
	err := c.BindJSON(&pond)
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}

	if pond.FarmID == "" || pond.ID == "" || pond.Name == "" {
		c.JSON(500, gin.H{
			"status":  500,
			"message": "please read documentation, system need more data",
		})
		return
	}

	var isDuplicateEntry bool
	// memanggil usecase untuk membuat item baru
	_, isDuplicateEntry, err = h.pondUsecase.Store(&pond)
	if isDuplicateEntry {
		c.JSON(409, gin.H{
			"status":  409,
			"message": "data pond '" + pond.ID + "' pada farm '" + pond.FarmID + "' sudah pernah dibuat",
		})
		return
	}
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": "input data berhasil",
	})
}

// UpdateById memperbarui Pond berdasarkan farm_id
func (h *PondHandler) UpdateById(c *gin.Context) {
	// mengambil body JSON farm dari Body Post
	var pond model.Pond
	err := c.BindJSON(&pond)
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}

	// memanggil usecase untuk memperbarui item
	isCreateData, err := h.pondUsecase.UpdateById(&pond)
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	if isCreateData {
		c.JSON(201, gin.H{
			"status":  201,
			"message": "pembuatan data '" + pond.ID + " pada farm '" + pond.FarmID + "' berhasil",
		})
		return
	}
	c.JSON(202, gin.H{
		"status":  202,
		"message": "pembaharuan data '" + pond.ID + "' pada farm '" + pond.FarmID + "' berhasil",
	})
}

// SoftDeleteById menghapus data Pond berdasarkan farm_id
func (h *PondHandler) SoftDeleteById(c *gin.Context) {
	farmId := c.Param("farm_id")
	pondId := c.Param("pond_id")
	isDataNotFound, err := h.pondUsecase.SoftDeleteById(farmId, pondId)
	if isDataNotFound {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "pond data '" + pondId + "' pada farm '" + farmId + "' tidak ada",
		})
		return
	}
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": "hapus data berhasil",
	})
}
