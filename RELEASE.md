## 0.0.1

1. 实现基本的功能，多个端口通过 http 请求映射到多个端口

## 0.0.2

1. 支持配置token

## 0.0.3

1. 优化代码, WebSocket 抽象为 Conn
2. 增加多路复用， 配置 smux = "true"

## 0.0.4

1. 优化错误信息

## 0.0.5

1. 重构代码, 修复 smux
2. 提供 `client` 配置 `mode` ，可选 `http` 和 `websocket`, 默认为 `websocket`
3. 提供 `client` 配置 `smux`，可选 `true` 和 `false`, 默认为 `false`