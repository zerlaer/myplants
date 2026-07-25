package router

import (
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"myplants/internal/config"
	"myplants/internal/controller"
)

// Setup 配置路由
func Setup(cfg *config.Config) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// 静态文件: 上传的图片
	r.Static("/uploads", cfg.Upload.Path)

	// API 路由组
	api := r.Group("/api")
	{
		// 仪表盘
		api.GET("/dashboard", controller.GetDashboard)

		// 植物
		plants := api.Group("/plants")
		{
			plants.GET("", controller.ListPlants)
			plants.GET("/:id", controller.GetPlant)
			plants.POST("", controller.CreatePlant)
			plants.PUT("/:id", controller.UpdatePlant)
			plants.DELETE("/:id", controller.DeletePlant)
			plants.GET("/:id/stats", controller.GetPlantStats)
		}

		// 养护提醒
		api.GET("/reminders", controller.GetReminders)

		// 养护记录
		care := api.Group("/care")
		{
			care.GET("", controller.ListCareRecords)
			care.POST("", controller.CreateCareRecord)
			care.DELETE("/:id", controller.DeleteCareRecord)
			care.POST("/one-click/:id/:type", controller.OneClickCare)
			care.POST("/batch", controller.BatchCare)
		}

		// 相册
		photos := api.Group("/photos")
		{
			photos.GET("", controller.ListPhotos)
			photos.POST("", controller.UploadPhoto)
			photos.POST("/avatar", controller.UploadAvatar)
			photos.PUT("/:id", controller.UpdatePhoto)
			photos.DELETE("/:id", controller.DeletePhoto)
		}

		// 花盆
		pots := api.Group("/pots")
		{
			pots.GET("", controller.ListPots)
			pots.POST("", controller.CreatePot)
			pots.PUT("/:id", controller.UpdatePot)
			pots.DELETE("/:id", controller.DeletePot)
		}

		// 换盆
		repot := api.Group("/repottings")
		{
			repot.GET("", controller.ListRepottings)
			repot.POST("", controller.CreateRepotting)
			repot.DELETE("/:id", controller.DeleteRepotting)
		}

		// 笔记
		notes := api.Group("/notes")
		{
			notes.GET("", controller.ListNotes)
			notes.POST("", controller.CreateNote)
			notes.DELETE("/:id", controller.DeleteNote)
		}
	}

	// 前端静态文件(构建后的 dist), 不存在时返回提示
	distPath := filepath.Join("frontend", "dist")
	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.File(filepath.Join(distPath, "assets", c.Param("filepath")))
	})
	r.NoRoute(func(c *gin.Context) {
		indexFile := filepath.Join(distPath, "index.html")
		c.File(indexFile)
	})

	return r
}
