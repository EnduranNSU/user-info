package httpin

import "github.com/gin-gonic/gin"

func NewGinRouter(h *UserInfoHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.StaticFile("/openapi.yaml", "docs/swagger.yaml")

	api := r.Group("/api/v1")
	{
		api.POST("/user-info", h.Create)
		api.GET("/user-info/latest", h.GetLatest)
		api.GET("/user-info", h.List)
	}

	return r
}
