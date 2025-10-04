package transaction

import "gorm.io/gorm"

// TransactionManagerFactory 事务管理器工厂
type TransactionManagerFactory struct {
	tm TransactionManager
}

// NewTransactionManagerFactory 事务管理器工厂构造函数
func NewTransactionManagerFactory(dbProvider func() *gorm.DB) *TransactionManagerFactory {
	// 使用 TransactionManager 的构造函数创建 TransactionManager
	tm := NewTransactionManager(dbProvider)

	return &TransactionManagerFactory{
		tm: tm,
	}
}

// TransactionManager 从事务管理器工厂获取事务管理器
func (f *TransactionManagerFactory) TransactionManager() TransactionManager {
	return f.tm
}
