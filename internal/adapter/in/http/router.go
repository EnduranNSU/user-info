package httpin

import (
	"github.com/gin-gonic/gin"

	_ "github.com/EnduranNSU/end-user-info/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGinRouter(h *UserInfoHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.StaticFile("/openapi.yaml", "docs/swagger.yaml")
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.POST("/user-info", h.Create)
		api.GET("/user-info/latest", h.GetLatest)
		api.GET("/user-info", h.List)
	}

	return r
}
