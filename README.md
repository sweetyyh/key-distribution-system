# key-distribution-system

基于《卡密销售平台技术方案 v2.0》的初始化工程骨架：

- Go + Gin 的 HTTP 服务入口
- 基础配置加载（Viper）
- 核心数据模型（GORM）
- 与方案一致的 API 路由占位
- 初始数据库建表脚本

## 快速启动

```bash
go mod tidy
go run ./cmd/server
```

默认配置文件：`configs/config.yaml`。
可通过环境变量 `KDS_CONFIG` 指定其他配置路径。

## 当前状态

当前仓库已经完成 P0/P1 的目录和接口骨架初始化，具体业务逻辑（用户、商品、订单、支付、卡密发放等）待继续实现。
