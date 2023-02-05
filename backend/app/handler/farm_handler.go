package handler

import (
	"aquafarm-management/app/model"
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
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

// Fetch mengambil semua item
func (h *FarmHandler) Fetch(c *gin.Context) {
	// memanggil usecase untuk mengambil item
	farms, err := h.farmUsecase.Fetch()
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

// GetById mengambil Farm berdasarkan ID
func (h *FarmHandler) GetById(c *gin.Context) {
	// mengambil ID item dari URL
	farmId := c.Param("farm_id")
	// memanggil usecase untuk mengambil item
	item, isNotFound, err := h.farmUsecase.GetById(farmId)
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
			"message": "farm data '" + farmId + "' tidak ada",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"data":   item,
	})
}

// Store membuat Farm baru
func (h *FarmHandler) Store(c *gin.Context) {
	// mengambil data item dari permintaan HTTP
	var farm model.Farm
	err := c.BindJSON(&farm)
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}

	if farm.ID == "" || farm.Name == "" {
		c.JSON(500, gin.H{
			"status":  500,
			"message": "please read documentation, system need more data",
		})
		return
	}

	var isDuplicateEntry bool
	// memanggil usecase untuk membuat item baru
	_, isDuplicateEntry, err = h.farmUsecase.Store(&farm)
	if isDuplicateEntry {
		c.JSON(409, gin.H{
			"status":  409,
			"message": "data farm dengan farm_id '" + farm.ID + "' sudah pernah dibuat",
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

// UpdateById memperbarui Farm berdasarkan farm_id
func (h *FarmHandler) UpdateById(c *gin.Context) {
	// mengambil body JSON farm dari Body Post
	var farm model.Farm
	err := c.BindJSON(&farm)
	if err != nil {
		// menangani error
		c.JSON(500, gin.H{
			"status":  500,
			"message": err.Error(),
		})
		return
	}

	// memanggil usecase untuk memperbarui item
	isCreateData, err := h.farmUsecase.UpdateById(&farm)
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
			"message": "pembuatan data '" + farm.ID + " berhasil",
		})
		return
	}
	c.JSON(202, gin.H{
		"status":  202,
		"message": "pembaharuan data '" + farm.ID + "' berhasil",
	})
}

// SoftDeleteById menghapus data Farm berdasarkan farm_id
func (h *FarmHandler) SoftDeleteById(c *gin.Context) {
	farmId := c.Param("farm_id")
	isDataNotFound, err := h.farmUsecase.SoftDeleteById(farmId)
	if isDataNotFound {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "farm data '" + farmId + "' tidak ada",
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
