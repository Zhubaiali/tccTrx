## 事务管理器（TXManager）流程

### 1. 结构定义

`TXManager`结构体定义了事务管理器的主要组件，包括：
- `ctx`和`stop`：用于控制事务管理器的生命周期。
- `opts`：事务管理器的配置选项。
- `txStore`：事务日志存储模块，用于存储事务的状态和日志。
- `registryCenter`：TCC组件注册中心，用于管理TCC组件。

### 2. 创建事务管理器

`NewTXManager`函数用于创建一个新的事务管理器，并启动其运行循环。

### 3. 事务处理

`Transaction`方法用于处理一个新的事务请求。它首先获取所有的TCC组件，然后创建事务明细记录，并执行两阶段提交（try-confirm/cancel）。

### 4. 两阶段提交

`twoPhaseCommit`方法实现了两阶段提交的逻辑。在第一阶段（try阶段），它并发执行所有TCC组件的try操作。如果所有try操作都成功，则进入第二阶段并执行confirm操作；否则执行cancel操作。

### 5. 事务监控

`run`方法是事务管理器的主运行循环。它定期检查hanging状态的事务，并尝试推进它们的状态。如果出现失败，它会使用back-off策略来避免重复执行。

### 6. 状态推进

`batchAdvanceProgress`和`advanceProgress`方法用于推进事务的状态。它们对每笔事务进行状态推进，根据事务的状态执行confirm或cancel操作，并更新事务的状态。

### 7. 锁定和解锁

在`run`方法中，通过调用`txStore.Lock`和`txStore.Unlock`来加锁和解锁事务日志存储模块，避免监控任务的重复执行。

### 8. 获取TCC组件

`getComponents`方法用于获取事务请求中的所有TCC组件，并检查它们的合法性

## 单测演示
![image](https://github.com/Zhubaiali/tccTrx/assets/69970253/33f7570e-b586-4bcf-83f5-ab92feb523ef)

