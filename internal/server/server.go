// Package server
//
// @title Ech0 API 文档
// @version 1.0
// @description 开源、自托管轻量级发布平台 Ech0 的 API 文档
// @host localhost:6277
// @BasePath /api
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lin-snow/ech0/internal/cache"
	"github.com/lin-snow/ech0/internal/config"
	"github.com/lin-snow/ech0/internal/database"
	"github.com/lin-snow/ech0/internal/di"
	"github.com/lin-snow/ech0/internal/event"
	commonModel "github.com/lin-snow/ech0/internal/model/common"
	"github.com/lin-snow/ech0/internal/router"
	"github.com/lin-snow/ech0/internal/task"
	"github.com/lin-snow/ech0/internal/transaction"
	errUtil "github.com/lin-snow/ech0/internal/util/err"
)

// Server 服务器结构体，包含Gin引擎
type Server struct {
	GinEngine      *gin.Engine
	httpServer     *http.Server          // 用于优雅停止服务器
	tasker         *task.Tasker          // 任务器
	eventRegistrar *event.EventRegistrar // 事件注册器
}

// New 创建一个新的服务器实例
func New() *Server {
	return &Server{}
}

// Init 初始化服务器
func (s *Server) Init() {
	// Mode
	if config.Config.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Gin Engine
	s.GinEngine = gin.New()

	// Database
	database.InitDatabase()

	// CacheFactory
	cacheFactory := cache.NewCacheFactory()

	// TransactionManagerFactory
	transactionManagerFactory := transaction.NewTransactionManagerFactory(database.GetDB)

	// Event System
	event.InitEventBus()

	// Handlers
	handlers, err := di.BuildHandlers(database.GetDB, cacheFactory, transactionManagerFactory, event.GetEventBus)
	if err != nil {
		errUtil.HandlePanicError(&commonModel.ServerError{
			Msg: commonModel.INIT_HANDLERS_PANIC,
			Err: err,
		})
	}

	// Router
	router.SetupRouter(s.GinEngine, handlers)

	// Tasker
	s.tasker, err = di.BuildTasker(database.GetDB, cacheFactory, transactionManagerFactory, event.GetEventBus)
	if err != nil {
		errUtil.HandlePanicError(&commonModel.ServerError{
			Msg: commonModel.INIT_TASKER_PANIC,
			Err: err,
		})
	}

	// EventRegistrar
	s.eventRegistrar, err = di.BuildEventRegistrar(
		database.GetDB,
		event.GetEventBus,
		cacheFactory,
		transactionManagerFactory,
	)
}

// Start 异步启动服务器
func (s *Server) Start() {
	port := config.Config.Server.Port
	PrintGreetings(port)

	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: s.GinEngine,
	}

	// 启动服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errUtil.HandlePanicError(&commonModel.ServerError{
				Msg: commonModel.GIN_RUN_FAILED,
				Err: err,
			})
		}
	}()
	// log.Println("🚀 Ech0 Server已启动")

	// 启动任务器
	go s.tasker.Start()
	// log.Println("🤖 任务器已启动")

	// 注册事件
	go s.eventRegistrar.Register()
	// log.Println("🔮 事件注册器已启动")
}

// Stop 优雅停止服务器
func (s *Server) Stop(ctx context.Context) error {
	// 使用传入的 context，如果没有则创建默认的 5 秒超时
	shutdownCtx := ctx
	var cancel context.CancelFunc

	if ctx == nil {
		shutdownCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	if s.httpServer == nil {
		fmt.Println("⚠️ HTTP 服务器未启动，无需关闭")
		return nil
	}

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}

	// 停止任务器
	s.tasker.Stop()

	// 等待事件系统任务结束
	s.eventRegistrar.Wait()

	return nil
}
