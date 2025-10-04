//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/lin-snow/ech0/internal/cache"
	backupHandler "github.com/lin-snow/ech0/internal/handler/backup"
	commonHandler "github.com/lin-snow/ech0/internal/handler/common"
	connectHandler "github.com/lin-snow/ech0/internal/handler/connect"
	echoHandler "github.com/lin-snow/ech0/internal/handler/echo"
	fediverseHandler "github.com/lin-snow/ech0/internal/handler/fediverse"
	settingHandler "github.com/lin-snow/ech0/internal/handler/setting"
	todoHandler "github.com/lin-snow/ech0/internal/handler/todo"
	userHandler "github.com/lin-snow/ech0/internal/handler/user"
	webHandler "github.com/lin-snow/ech0/internal/handler/web"
	commonRepository "github.com/lin-snow/ech0/internal/repository/common"
	connectRepository "github.com/lin-snow/ech0/internal/repository/connect"
	echoRepository "github.com/lin-snow/ech0/internal/repository/echo"
	fediverseRepository "github.com/lin-snow/ech0/internal/repository/fediverse"
	keyvalueRepository "github.com/lin-snow/ech0/internal/repository/keyvalue"
	todoRepository "github.com/lin-snow/ech0/internal/repository/todo"
	userRepository "github.com/lin-snow/ech0/internal/repository/user"
	backupService "github.com/lin-snow/ech0/internal/service/backup"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	connectService "github.com/lin-snow/ech0/internal/service/connect"
	echoService "github.com/lin-snow/ech0/internal/service/echo"
	fediverseService "github.com/lin-snow/ech0/internal/service/fediverse"
	settingService "github.com/lin-snow/ech0/internal/service/setting"
	todoService "github.com/lin-snow/ech0/internal/service/todo"
	userService "github.com/lin-snow/ech0/internal/service/user"
	"github.com/lin-snow/ech0/internal/task"
	"github.com/lin-snow/ech0/internal/transaction"
	"gorm.io/gorm"
)

// BuildHandlers 使用wire生成的代码来构建Handlers实例
func BuildHandlers(
	dbProvider func() *gorm.DB,
	cacheFactory *cache.CacheFactory,
	tmFactory *transaction.TransactionManagerFactory,
) (*Handlers, error) {
	wire.Build(
		CacheSet,
		TransactionManagerSet,
		WebSet,
		UserSet,
		EchoSet,
		CommonSet,
		SettingSet,
		TodoSet,
		ConnectSet,
		BackupSet,
		FediverseSet,
		NewHandlers, // NewHandlers 聚合各个模块的Handler
	)

	return &Handlers{}, nil
}

func BuildTasker(
	dbProvider func() *gorm.DB,
	cacheFactory *cache.CacheFactory,
	tmFactory *transaction.TransactionManagerFactory,
) (*task.Tasker, error) {
	wire.Build(
		CacheSet,
		TransactionManagerSet,
		EchoSet,
		CommonSet,
		SettingSet,
		TaskSet,
	)
	return &task.Tasker{}, nil
}

// CacheSet 包含了构建缓存所需的所有 Provider
var CacheSet = wire.NewSet(
	ProvideCache,
)

// TransactionManagerSet 包含了构建事务管理器所需的所有 Provider
var TransactionManagerSet = wire.NewSet(
	ProvideTransactionManager,
)

// WebSet 包含了构建 WebHandler 所需的所有 Provider
var WebSet = wire.NewSet(
	webHandler.NewWebHandler,
)

// UserSet 包含了构建 UserHandler 所需的所有 Provider
var UserSet = wire.NewSet(
	userRepository.NewUserRepository,
	userService.NewUserService,
	userHandler.NewUserHandler,
)

// EchoSet 包含了构建 EchoHandler 所需的所有 Provider
var EchoSet = wire.NewSet(
	echoRepository.NewEchoRepository,
	echoService.NewEchoService,
	echoHandler.NewEchoHandler,
)

// CommonSet 包含了构建 CommonHandler 所需的所有 Provider
var CommonSet = wire.NewSet(
	commonRepository.NewCommonRepository,
	commonService.NewCommonService,
	commonHandler.NewCommonHandler,
)

// SettingSet 包含了构建 SettingHandler 所需的所有 Provider
var SettingSet = wire.NewSet(
	keyvalueRepository.NewKeyValueRepository,
	settingService.NewSettingService,
	settingHandler.NewSettingHandler,
)

// TodoSet 包含了构建 TodoHandler 所需的所有 Provider
var TodoSet = wire.NewSet(
	todoRepository.NewTodoRepository,
	todoService.NewTodoService,
	todoHandler.NewTodoHandler,
)

// ConnectSet 包含了构建 ConnectHandler 所需的所有 Provider
var ConnectSet = wire.NewSet(
	connectRepository.NewConnectRepository,
	connectService.NewConnectService,
	connectHandler.NewConnectHandler,
)

// BackupSet 包含了构建 BackupHandler 所需的所有 Provider
var BackupSet = wire.NewSet(
	backupHandler.NewBackupHandler,
	backupService.NewBackupService,
)

// TaskSet 包含了构建 Tasker 所需的所有 Provider
var TaskSet = wire.NewSet(
	task.NewTasker,
)

// FediverseSet 包含了构建 Fediverse 所需的所有 Provider
var FediverseSet = wire.NewSet(
	fediverseRepository.NewFediverseRepository,
	fediverseService.NewFediverseService,
	fediverseHandler.NewFediverseHandler,
)
