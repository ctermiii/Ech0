package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lin-snow/ech0/internal/middleware"
)

// setupMiddleware 设置中间件
func setupMiddleware(r *gin.Engine) {
	// Cors middleware
	r.Use(middleware.Cors())
	// Global write guard middleware
	r.Use(middleware.WriteGuard())
}
