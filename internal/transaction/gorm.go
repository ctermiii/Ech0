package transaction

import (
	"context"

	"gorm.io/gorm"
)

// GormTransactionManager 实现了 TransactionManager 接口，使用 GORM 进行事务管理
type GormTransactionManager struct {
	dbProvider func() *gorm.DB
}

func NewGormTransactionManager(dbProvider func() *gorm.DB) *GormTransactionManager {
	return &GormTransactionManager{
		dbProvider: dbProvider,
	}
}

// Run 在 GormTransactionManager 中实现了 TransactionManager 接口的 Run 方法
// 该方法接受一个函数 fn，并在一个新的事务中执行它
// 函数 fn 接受一个 context.Context 参数，表示当前事务的上下文
// 如果 fn 执行成功，事务将被提交；如果 fn 返回错误，事务将被回滚
// 参数:
//   - fn: 一个函数，接受一个 context.Context 参数，并返回一个 error
//
// 返回:
//   - error: 如果 fn 执行成功返回 nil，否则返回错误信息
func (tm *GormTransactionManager) Run(fn func(ctx context.Context) error) error {
	// 返回一个新的事务上下文
	// 在这个上下文中，txKey 被设置为当前事务的 gorm.DB，使用gorm自带的自动事务管理
	return tm.dbProvider().Transaction(func(tx *gorm.DB) error {
		// 将当前事务的 gorm.DB 设置到上下文中，这里创建一个新的上下文
		ctx := context.WithValue(context.Background(), TxKey, tx)

		// 执行传入的函数，并传递事务上下文
		return fn(ctx)
	})
}
