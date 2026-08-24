# Bug Reproduction

## Bug

HTTP timeout middleware、遥测与事件 handler，以及订阅 Create/List 把请求 context 替换成 `context.Background()`，取消信号无法到达后端。

## Trigger

取消请求或触发请求 deadline 后调用订阅、事件和遥测入口。慢速仓储仍持续执行，取消 context 的分支不会触发。

## Error

红测中四类场景均继续执行约 301ms，慢仓储记录 `ctx NOT canceled`，而预期应在取消后立即返回。
