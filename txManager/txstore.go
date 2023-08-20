package txManager

import (
	"context"
	"tccTrx/component"
	"time"
)

// TXStore 事务日志管理
type TXStore interface {
	// CreateTX 创建一条事务
	CreateTX(ctx context.Context, components ...component.TCCComponent) (txID string, err error)
	// TXUpdate 更新事务进度：
	// 规则为：倘若有一个 component try 操作执行失败，则整个事务失败；倘若所有 component try 操作执行成功，则事务成功
	TXUpdate(ctx context.Context, txID string, componentID string, accept bool) error
	// GetHangingTXs 获取到所有处于中间态的事务
	GetHangingTXs(ctx context.Context) ([]*Transaction, error)
	// Lock 锁住事务日志表
	Lock(ctx context.Context, expireDuration time.Duration) error
	// Unlock 解锁事务日志表
	Unlock(ctx context.Context) error
}
