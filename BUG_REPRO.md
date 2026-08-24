# Bug Reproduction

## Bug

Health、PostgreSQL 探活以及迁移 Up/Down 在调用方 context 已取消时仍使用独立的后台 context，因而继续执行并返回成功。

## Trigger

向这四个存储入口传入已经取消的 context。探活返回 nil；迁移操作继续处理迁移列表并改变 Applied 状态。

## Error

红测中四个定向测试均失败：`Ping() error = <nil>, want context.Canceled`、`Up() error = <nil>, want context.Canceled`、`Down() error = <nil>, want context.Canceled`。
