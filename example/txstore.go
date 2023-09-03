package example

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/demdxx/gocast"
	"github.com/xiaoxuxiansheng/redis_lock"
	"tccTrx/component"
	"tccTrx/example/dao"
	"tccTrx/example/pkg"
	"tccTrx/txManager"
	"time"
)

type MockTXStore struct {
	client *redis_lock.Client
	dao    *dao.TXRecordDAO
}

func NewMockTXStore(dao *dao.TXRecordDAO, client *redis_lock.Client) *MockTXStore {
	return &MockTXStore{
		dao:    dao,
		client: client,
	}
}

func (m *MockTXStore) CreateTX(ctx context.Context, components ...component.TCCComponent) (string, error) {
	// 创建一项内容，里面以唯一事务 id 为 key
	componentTryStatuses := make(map[string]*dao.ComponentTryStatus, len(components))
	for _, component := range components {
		componentTryStatuses[component.ID()] = &dao.ComponentTryStatus{
			ComponentID: component.ID(),
			TryStatus:   txManager.TryHanging.String(),
		}
	}

	statusesBody, _ := json.Marshal(componentTryStatuses)
	txID, err := m.dao.CreateTXRecord(ctx, &dao.TXRecordPO{
		Status:               txManager.TXHanging.String(),
		ComponentTryStatuses: string(statusesBody),
	})
	if err != nil {
		return "", err
	}

	return gocast.ToString(txID), nil
}

func (m *MockTXStore) TXUpdate(ctx context.Context, txID string, componentID string, accept bool) error {
	_txID := gocast.ToUint(txID)
	status := txManager.TXFailure.String()
	if accept {
		status = txManager.TXSuccessful.String()
	}
	return m.dao.UpdateComponentStatus(ctx, _txID, componentID, status)
}

func (m *MockTXStore) GetHangingTXs(ctx context.Context) ([]*txManager.Transaction, error) {
	records, err := m.dao.GetTXRecords(ctx, dao.WithStatus(txManager.TryHanging))
	if err != nil {
		return nil, err
	}

	txs := make([]*txManager.Transaction, 0, len(records))
	for _, record := range records {
		componentTryStatuses := make(map[string]*dao.ComponentTryStatus)
		_ = json.Unmarshal([]byte(record.ComponentTryStatuses), &componentTryStatuses)
		components := make([]*txManager.ComponentTryEntity, 0, len(componentTryStatuses))
		for _, component := range componentTryStatuses {
			components = append(components, &txManager.ComponentTryEntity{
				ComponentID: component.ComponentID,
				TryStatus:   txManager.ComponentTryStatus(component.TryStatus),
			})
		}

		txs = append(txs, &txManager.Transaction{
			TXID:       gocast.ToString(record.ID),
			Status:     txManager.TXHanging,
			CreatedAt:  record.CreatedAt,
			Components: components,
		})
	}

	return txs, nil
}

func (m *MockTXStore) Lock(ctx context.Context, expireDuration time.Duration) error {
	lock := redis_lock.NewRedisLock(pkg.BuildTXRecordLockKey(), m.client, redis_lock.WithExpireSeconds(int64(expireDuration.Seconds())))
	return lock.Lock(ctx)
}

func (m *MockTXStore) Unlock(ctx context.Context) error {
	lock := redis_lock.NewRedisLock(pkg.BuildTXRecordLockKey(), m.client)
	return lock.Unlock(ctx)
}

// 提交事务的最终状态
func (m *MockTXStore) TXSubmit(ctx context.Context, txID string, success bool) error {
	do := func(ctx context.Context, dao *dao.TXRecordDAO, record *dao.TXRecordPO) error {
		if success {
			record.Status = txManager.TXSuccessful.String()
		} else {
			record.Status = txManager.TXFailure.String()
		}
		return dao.UpdateTXRecord(ctx, record)
	}
	return m.dao.LockAndDo(ctx, gocast.ToUint(txID), do)
}

// 获取指定的一笔事务
func (m *MockTXStore) GetTX(ctx context.Context, txID string) (*txManager.Transaction, error) {
	records, err := m.dao.GetTXRecords(ctx, dao.WithID(gocast.ToUint(txID)))
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, errors.New("get tx failed")
	}

	componentTryStatuses := make(map[string]*dao.ComponentTryStatus)
	_ = json.Unmarshal([]byte(records[0].ComponentTryStatuses), &componentTryStatuses)

	components := make([]*txManager.ComponentTryEntity, 0, len(componentTryStatuses))
	for _, tryItem := range componentTryStatuses {
		components = append(components, &txManager.ComponentTryEntity{
			ComponentID: tryItem.ComponentID,
			TryStatus:   txManager.ComponentTryStatus(tryItem.TryStatus),
		})
	}
	return &txManager.Transaction{
		TXID:       txID,
		Status:     txManager.TXStatus(records[0].Status),
		Components: components,
		CreatedAt:  records[0].CreatedAt,
	}, nil
}
