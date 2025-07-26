package service

import "github.com/gin-gonic/gin"

type BackupServiceInterface interface {
	// 执行备份
	Backup(userid uint) error

	// 导出备份
	ExportBackup(userid uint, ctx *gin.Context) error

	// 恢复备份
	// ImportBackup(userid uint) error
}
