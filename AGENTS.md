# AGENTS.md
- This is the only AGENTS.md, there are no recursive AGENTS.md

## 项目概述
Cascade-oj 是一个在线代码测评系统，采用微服务架构
- 后端：Go 1.25 + Kratos v2.9，多微服务；使用 MySQL、Redis、RabbitMQ；测评机基于 go-judge。
- 前端：Vue 3 + TypeScript，pnpm monorepo 组织子包。
- 部署与测试：Docker Compose（专注于容器化系统测试，无热重载）
- 更多具体说明参考 `README.md`

## 快速指令
### 后端
后端目录 `cascade-oj/`
所有后端操作均需在 Docker 容器中执行，无本地热重载
- 构建但不启动单个或多个服务镜像：`docker compose -f docker-compose.dependencies.yml -f docker-compose.yml build <service1> [service2...]`
- 构建并启动单个或多个服务镜像：`./deploy.sh --build <service1> [service2...]`
- 启动服务（已构建镜像）：`./deploy.sh <service1> [service2...]`
- 查看其他后端辅助指令：`make help`（参考 Makefile）

构建镜像示例：
```bash
./deploy.sh --build public user judge admin
```

启动示例：
```bash
./deploy.sh public user judge admin
```

### 前端
前端目录 `cascade-oj-frontend/`
- 安装依赖（锁定版本）：`pnpm install --frozen-lockfile`
- 构建单个 monorepo 子包（前端产物编译）：`pnpm run <project_name>:build`
- 构建并启动所有分发与网关容器：`./setup.sh --build`（需先完成前端产物编译）
- 启动完整前端+网关：`./setup.sh`

`<project_name>` 为 monorepo 子包名称，包含 login、competition、admin

## 开发测试流程
`修改 -> 构建 -> 启动 -> 人类用户验证`
使用 Docker Compose 启动全部依赖和后端、前端服务，模拟生产环境。

## 代码规范
### 后端（Golang）
- 命名：遵循官方标准（驼峰式，导出标识符首字母大写）
- 错误处理：必须显式处理 error，忽略需用 `_` 并注释原因
- 单元测试：使用 `go test`，集成测试使用 Docker Compose，并交给人类用户验证
- 依赖管理：禁止手动更新依赖，增加或更新依赖需要经过人类用户同意。
其他细节与原代码保持一致，支持 `go fmt`。
- wire 胶水代码生成：go mod tidy 可能将 wire 的依赖从 go.mod 删除，使用 `go get github.com/google/wire/cmd/wire@v0.7.0` 补全

### 前端（Vue3 + TypeScript）
- 语言：TypeScript 禁止使用 `any`。Vue3 使用组合式 API (Composition API) 风格
- 组件：使用 `<script setup lang="ts">`
- 样式：使用独立 CSS 文件，scoped CSS 导入，避免嵌入 vue 文件
- 依赖管理：严格使用 `pnpm install --frozen-lockfile`，增加或更新依赖需要经过人类用户同意。
其他细节与原代码保持一致。

## 操作边界与禁止行为
禁止提交：测试配置文件、硬编码测试数据（如测试 token、硬编码的用户凭据）。
换行与编码：LF 换行，UTF-8，无 BOM。
无热重载：任何代码更改后，需重新构建容器镜像并重启服务。

## 调试与排错
### 单元测试
Go 内置 `testing` 包，在存在单元测试的包目录运行 `go test ./...`
### 集成测试
docker-compose 环境，使用命令启动需要的容器，人类用户执行测试
### 系统测试
涉及完整业务链条（如从提交代码到判题返回）的测试，必须使用 Docker Compose 启动全部服务（前端网关 + 后端微服务 + 基础设施）

### 错误排查
先看静态类型检查结果，再提醒人类用户依次检查以下环节：
- 前端请求响应
- 后端日志收到请求
- 后端处理结果
前端请求响应在测试浏览器打开开发工具（F12），检查控制台输出、网络请求
后端部分排查可以检查日志文件（日志文件极大，禁止阅读，必须让人类用户手动检查）
日志文件位于 `cascade-oj/data/logs/`
