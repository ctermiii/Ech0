package database

import (
	"errors"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/lin-snow/ech0/internal/config"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	connectModel "github.com/lin-snow/ech0/internal/model/connect"
	echoModel "github.com/lin-snow/ech0/internal/model/echo"
	fediverseModel "github.com/lin-snow/ech0/internal/model/fediverse"
	todoModel "github.com/lin-snow/ech0/internal/model/todo"
	userModel "github.com/lin-snow/ech0/internal/model/user"

	util "github.com/lin-snow/ech0/internal/util/err"
	"gorm.io/driver/sqlite"

	// "github.com/glebarez/sqlite" // 使用 glebarez/sqlite 作为 SQLite 驱动
	"gorm.io/gorm"
)

// DB 全局数据库连接变量
// var DB *gorm.DB

// 使用 atomic.Value 来存储 *gorm.DB，确保线程安全和支持热更新
var db atomic.Value // 用于存储 *gorm.DB

var writeLocked atomic.Bool

func GetDB() *gorm.DB {
	return db.Load().(*gorm.DB)
}

func SetDB(newDB *gorm.DB) {
	db.Store(newDB)
}

func DBProvider() func() *gorm.DB {
	return GetDB
}

// EnableWriteLock 启用写锁，阻止新的写操作
func EnableWriteLock() {
	writeLocked.Store(true)
}

// DisableWriteLock 关闭写锁，允许写操作
func DisableWriteLock() {
	writeLocked.Store(false)
}

// SetWriteLock 手动设置写锁状态
func SetWriteLock(enabled bool) {
	writeLocked.Store(enabled)
}

// IsWriteLocked 判断当前是否启用了写锁
func IsWriteLocked() bool {
	return writeLocked.Load()
}

// InitDatabase 初始化数据库连接
func InitDatabase() {
	// 读取数据库类型和保存路径
	dbType := config.Config.Database.Type
	dbPath := config.Config.Database.Path

	dir := dbPath[:len(dbPath)-len("/ech0.db")] // 提取目录部分
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		util.HandlePanicError(&commonModel.ServerError{
			Msg: commonModel.CREATE_DB_PATH_PANIC,
			Err: err,
		})
	}

	if dbType == "sqlite" {
		// 添加 PRAGMA 参数，例如 WAL 模式和外键支持
		// pragma := config.Config.Database.Pragma // 从配置读取
		// dsn := dbPath + "?" + pragma
		var err error
		SQLiteDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			util.HandlePanicError(&commonModel.ServerError{
				Msg: commonModel.INIT_DATABASE_PANIC,
				Err: err,
			})
		}
		SetDB(SQLiteDB)
	}

	if err := MigrateDB(); err != nil {
		util.HandlePanicError(&commonModel.ServerError{
			Msg: commonModel.MIGRATE_DB_PANIC,
			Err: err,
		})
	}

	// 从 1.x 迁移到 2.x
	UpdateMigration()
}

// MigrateDB 执行数据库迁移
func MigrateDB() error {
	models := []interface{}{
		&userModel.User{},
		&echoModel.Echo{},
		&echoModel.Image{},
		&commonModel.KeyValue{},
		&todoModel.Todo{},
		&connectModel.Connected{},
		&commonModel.TempFile{},

		// Fediverse 相关
		&fediverseModel.Follow{},
		&fediverseModel.Follower{},
	}

	return GetDB().AutoMigrate(
		models...,
	)
}

// HotChangeDatabase 热切换数据库连接
func HotChangeDatabase(newDBPath string) error {
	// 获取当前数据库连接
	oldDB := GetDB()

	// 彻底关闭旧连接
	if oldDB != nil {
		if err := CloseDatabaseFully(oldDB); err != nil {
			return err
		}
	}

	// 打开新连接
	newDB, err := gorm.Open(sqlite.Open(newDBPath), &gorm.Config{})
	if err != nil {
		return err
	}

	SetDB(newDB)
	return nil
}

// CloseDatabaseFully 彻底关闭数据库连接，释放资源
func CloseDatabaseFully(db *gorm.DB) error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
		SetDB(nil)

		// 强制 GC 回收
		runtime.GC()
		time.Sleep(100 * time.Millisecond)

		return nil
	}

	return errors.New(commonModel.DATABASE_CLOSE_FAILED)
}
