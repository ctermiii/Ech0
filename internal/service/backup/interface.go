package service

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type BackupServiceInterface interface {
	// Backup 执行备份
	Backup(userid uint) error

	// ExportBackup 导出备份
	ExportBackup(ctx *gin.Context) error

	// 恢复备份
	ImportBackup(ctx *gin.Context, userid uint, file *multipart.FileHeader) error
}
