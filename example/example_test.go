package example

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"tccTrx/example/dao"
	"tccTrx/example/pkg"
	"tccTrx/txManager"
	"testing"
	"time"
)

const (
	dsn      = "debian-sys-maint:cCxZbowDDn7xWMzg@tcp(127.0.0.1:3306)/tcc"
	network  = "tcp"
	address  = "127.0.0.1:6379"
	password = ""
)

var rdb *redis.Client

func Test_Redis(t *testing.T) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// 验证 redis 连接是否正常，get一个存在的key "kk1":
	val2, _ := rdb.Get("kk1").Result()

	zsetKey := "language_rank"

	num, _ := rdb.ZAdd(zsetKey, redis.Z{Score: 100.0, Member: "C#"}).Result()
	fmt.Println(num)
	fmt.Println(val2)
}

func Test_TCC(t *testing.T) {
	redisClient := pkg.NewRedisClient(network, address, password)
	mysqlDB, err := pkg.NewDB(dsn)
	if err != nil {
		t.Error(err)
		return
	}

	componentAID := "componentA"
	componentBID := "componentB"
	componentCID := "componentC"

	// 构造出对应的 tcc component
	componentA := NewMockComponent(componentAID, redisClient)
	componentB := NewMockComponent(componentBID, redisClient)
	componentC := NewMockComponent(componentCID, redisClient)

	// 构造出事务日志存储模块
	txRecordDAO := dao.NewTXRecordDAO(mysqlDB)
	txStore := NewMockTXStore(txRecordDAO, redisClient)

	txmanager := txManager.NewTXManager(txStore, txManager.WithMonitorTick(time.Second))
	defer txmanager.Stop()

	// 完成各组件的注册
	if err := txmanager.Register(componentA); err != nil {
		t.Error(err)
		return
	}

	if err := txmanager.Register(componentB); err != nil {
		t.Error(err)
		return
	}

	if err := txmanager.Register(componentC); err != nil {
		t.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	success, err := txmanager.Transaction(ctx, []*txManager.RequestEntity{
		{ComponentID: componentAID,
			Request: map[string]interface{}{
				"biz_id": componentAID + "_biz",
			},
		},
		{ComponentID: componentBID,
			Request: map[string]interface{}{
				"biz_id": componentBID + "_biz",
			},
		},
		{ComponentID: componentCID,
			Request: map[string]interface{}{
				"biz_id": componentCID + "_biz",
			},
		},
	}...)
	if err != nil {
		t.Errorf("tx failed, err: %v", err)
		return
	}
	if !success {
		t.Error("tx failed")
		return
	}

	<-time.After(2 * time.Second)

	t.Log("success")
}
