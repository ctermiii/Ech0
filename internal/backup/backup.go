package backup

import (
	"errors"
	"fmt"
	"time"

	fileUtil "github.com/lin-snow/ech0/internal/util/file"
)

const (
	dataDir        = "data"                // 待备份的数据目录
	backupDir      = "backup"              // 备份后存储zip的目录
	backupFileName = "ech0_backup"         // 备份文件名
	excludeFile    = "*.log"               // 排除的文件名
	timeLayout     = "2006-01-02_15-04-05" // 时间格式化布局
)

// ExecuteBackup 执行备份
func ExecuteBackup() (string, string, error) {
	backupTime := time.Now().Format(timeLayout)
	backupFileName := fmt.Sprintf("%s_%s.zip", backupFileName, backupTime) // 暂时不开启多备份，每次只保留最新的一份备份
	backupPath := fmt.Sprintf("%s/%s", backupDir, backupFileName)

	return backupPath, backupFileName, fileUtil.ZipDirectoryWithOptions(dataDir, backupPath, fileUtil.ZipOptions{
		ExcludePatterns: []string{excludeFile},
	})
}

// ExecuteRestore 执行恢复
func ExecuteRestore(backupFilePath string) error {
	// 检查备份文件是否存在
	if !fileUtil.FileExists(backupFilePath) {
		return errors.New("备份文件不存在: " + backupFilePath)
	}

	// 解压备份文件到数据目录
	if err := fileUtil.UnzipFile(backupFilePath, dataDir); err != nil {
		return err
	}

	return nil
}
