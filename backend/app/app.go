package app

import (
	"aquafarm-management/app/config"
	"aquafarm-management/app/handler"
	"aquafarm-management/app/repository"
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type App struct {
	Router *gin.Engine
}

func New(cfg *config.Config) *App {
	a := &App{
		Router: gin.Default(),
	}

	db := repository.NewDB(cfg)

	repository.SetupDB(db)

	// INIT - Item
	newRepository := repository.NewItemRepository(db)
	itemUsecase := usecase.NewItemUsecase(newRepository)
	itemHandler := handler.NewItemHandler(itemUsecase)

	// INIT - Farm
	farmRepo := repository.NewFarmRepository(db)
	farmCase := usecase.NewFarmUsecase(farmRepo)
	farmHand := handler.NewFarmHandler(farmCase)

	// INIT - Pond
	pondRepo := repository.NewPondRepository(db)
	pondCase := usecase.NewPondUsecase(pondRepo, farmRepo)
	pondHand := handler.NewPondHandler(pondCase)

	// INIT - API
	authRepo := repository.NewAuthRepository(db)

	authorized := a.Router.Group("/")
	authorized.Use(func(c *gin.Context) {
		headerAPIKey := c.GetHeader("Authorization")
		accessId, err := authRepo.GetApiKey(headerAPIKey)
		if err != nil || accessId == nil || *accessId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Unauthorized",
			})
			return
		}
	})

	v1 := authorized.Group("/v1")
	{
		farm := v1.Group("/farm")
		{
			farm.GET("/", farmHand.Fetch)
			farm.GET("/:farm_id", farmHand.GetById)
			farm.POST("/", farmHand.Store)
			farm.PUT("/", farmHand.UpdateById)
			farm.DELETE("/:farm_id", farmHand.SoftDeleteById)
		}

		pond := v1.Group("/pond")
		{
			pond.GET("/", pondHand.Fetch)
			pond.GET("/:farm_id/:pond_id", pondHand.GetById)
			pond.POST("/", pondHand.Store)
			pond.PUT("/", pondHand.UpdateById)
			pond.DELETE("/:farm_id/:pond_id", pondHand.SoftDeleteById)
		}

		logs := v1.Group("/logs")
		{
			logs.GET("/", itemHandler.Fetch)
			logs.GET("/statistics", itemHandler.Fetch)
		}
	}

	return a
}

func (a *App) Run(cfg *config.Config) error {
	return a.Router.Run(":" + cfg.Server.Port)
}
