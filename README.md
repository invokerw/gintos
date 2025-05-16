# gintos

参考 kratos。使用 gin 封装一个简单的 web 框架，支持中间件、路由、参数解析、JSON 响应等功能。

## 协议生成
使用 protobuf grpc 定义接口协议。protoc-gen-go-gin 生成 gin 的路由和处理函数。

## gintos demo

gintos 实现例子

### 前端
https://github.com/invokerw/gintos-frontend

将前端代码放在 `assets/frontend` 目录下，使用 `go run .` 启动服务。
协议定义在 `api` 目录下，使用 `make api` 命令生成代码。

