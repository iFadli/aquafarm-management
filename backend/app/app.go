package app

import (
	"aquafarm-management/app/config"
	"aquafarm-management/app/handler"
	"aquafarm-management/app/model"
	"aquafarm-management/app/repository"
	"aquafarm-management/app/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strconv"
)

type App struct {
	Router *gin.Engine
}

func New(cfg *config.Config) *App {
	a := &App{
		Router: gin.Default(),
	}

	a.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	db := repository.NewDB(cfg)

	repository.SetupDB(db)

	// INIT - AUTH
	authRepo := repository.NewAuthRepository(db)
	//INIT - LOGS
	logRepo := repository.NewLogRepository(db)

	authorized := a.Router.Group("/")
	authorized.Use(func(c *gin.Context) {
		headerAPIKey := c.GetHeader("Authorization")
		accessId, err := authRepo.GetApiKey(headerAPIKey)
		if err != nil || accessId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			return
		}

		c.Next()
		// Buat log baru
		logger := &model.Logs{
			AccessID:  accessId,
			IpAddress: c.RemoteIP(),
			UserAgent: c.Request.UserAgent(),
			Request:   c.Request.Method + " " + c.Request.URL.Path,
			Response:  strconv.Itoa(c.Writer.Status()),
		}
		logRepo.FirstLog(logger)
	})

	// INIT - Farm
	farmRepo := repository.NewFarmRepository(db)
	farmCase := usecase.NewFarmUsecase(farmRepo)
	farmHand := handler.NewFarmHandler(farmCase)

	// INIT - Pond
	pondRepo := repository.NewPondRepository(db)
	pondCase := usecase.NewPondUsecase(pondRepo, farmRepo)
	pondHand := handler.NewPondHandler(pondCase)

	//INIT - Stat
	statCase := usecase.NewStatUsecase(logRepo)
	statHand := handler.NewStatHandler(statCase)

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
			logs.GET("/", statHand.FetchLogs)
			logs.GET("/statistics", statHand.Fetch)
		}
	}

	return a
}

func (a *App) Run(cfg *config.Config) error {
	return a.Router.Run(":" + cfg.Server.Port)
}
