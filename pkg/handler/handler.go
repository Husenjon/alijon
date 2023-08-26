package handler

import (
	"github.com/Husenjon/InkassBack/pkg/logger"
	"github.com/Husenjon/InkassBack/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services service.Service
	log      logger.Logger
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: *services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	auth := router.Group("/auth")
	{
		// auth.POST("sign-up", h.signUp)
		auth.POST("sign-in", h.signIn)
	}
	api := router.Group("api")
	{
		user := api.Group("/user", h.userIdentity)
		{
			user.POST("sign-up", h.signUp)
		}
		lists := api.Group("/lists")
		{
			lists.POST("/", func(ctx *gin.Context) {})
			lists.GET("/", func(ctx *gin.Context) {})
			lists.GET("/:id", func(ctx *gin.Context) {})
			lists.PUT("/:id", func(ctx *gin.Context) {})
			lists.DELETE("/:id", func(ctx *gin.Context) {})

			items := lists.Group(":id/items")
			{
				items.POST("/", func(ctx *gin.Context) {})
				items.GET("/", func(ctx *gin.Context) {})
				items.GET("/:item_id", func(ctx *gin.Context) {})
				items.PUT("/:item_id", func(ctx *gin.Context) {})
				items.DELETE("/:item_id", func(ctx *gin.Context) {})
			}
		}
	}
	return router
}
