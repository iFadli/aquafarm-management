package app

import (
	"aquafarm-management/app/config"
	"aquafarm-management/app/handler"
	"aquafarm-management/app/repository"
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
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

	v1 := a.Router.Group("/v1")
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
			pond.GET("/", itemHandler.Fetch)
			pond.GET("/:id", itemHandler.Get)
			pond.POST("/", itemHandler.Store)
			pond.PUT("/:id", itemHandler.Update)
			pond.DELETE("/:id", itemHandler.Delete)
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
