package handler

import "github.com/gin-gonic/gin"

type WebHandlerInterface interface {
	Templates() *gin.HandlerFunc
}
