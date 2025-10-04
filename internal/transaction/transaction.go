package transaction

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

const TxKey contextKey = "tx"

// TransactionManager 定义事务管理器接口
type TransactionManager interface {
	// Run 执行一个事务
	Run(fn func(ctx context.Context) error) error
}

func NewTransactionManager(dbProvider func() *gorm.DB) TransactionManager {
	// 使用GORM提供的事务管理器
	return NewGormTransactionManager(dbProvider)
}
