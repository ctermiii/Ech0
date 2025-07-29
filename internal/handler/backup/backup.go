package handler

import (
	"github.com/gin-gonic/gin"

	res "github.com/lin-snow/ech0/internal/handler/response"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	service "github.com/lin-snow/ech0/internal/service/backup"
)

type BackupHandler struct {
	backupService service.BackupServiceInterface
}

func NewBackupHandler(backupService service.BackupServiceInterface) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// Backup 执行备份
func (backupHandler *BackupHandler) Backup() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		userId := ctx.MustGet("userid").(uint)
		if err := backupHandler.backupService.Backup(userId); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.BACKUP_SUCCESS,
		}
	})
}

// ExportBackup 导出备份
func (backupHandler *BackupHandler) ExportBackup() gin.HandlerFunc {
	return res.Execute(func(ctx *gin.Context) res.Response {
		userId := ctx.MustGet("userid").(uint)
		if err := backupHandler.backupService.ExportBackup(userId, ctx); err != nil {
			return res.Response{
				Msg: "",
				Err: err,
			}
		}

		return res.Response{
			Msg: commonModel.EXPORT_BACKUP_SUCCESS,
		}
	})
}
