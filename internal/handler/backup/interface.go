package handler

import (
	"github.com/gin-gonic/gin"
)

type BackupHandlerInterface interface {
	// Backup 执行备份
	Backup() gin.HandlerFunc

	// ExportBackup 导出备份
	ExportBackup() gin.HandlerFunc

	// ImportBackup 恢复备份
	ImportBackup() gin.HandlerFunc
}
