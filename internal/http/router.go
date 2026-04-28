package http

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("/v1/check", h.Check)
	router.GET("/health", h.Health)

	router.NoRoute(func(c *gin.Context) {
		c.String(404, "404 page not found")
	})

	return router
}
