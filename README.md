<picture>
  <source media="(prefers-color-scheme: dark)" srcset="assets/cascade-dark.svg">
  <img alt="Logo" src="assets/cascade.svg">
</picture>

# Cascade-oj

用户友好的在线测评 (Online Judge) Web 平台，实训项目

项目后端使用 Kratos + golang 实现微服务，前端使用 Vue3 + typescript 开发，采用 pnpm monorepo 组织。

项目以竞赛、题库、提交评测为业务核心，使用 MySQL/Redis/RabbitMQ 基础设施，并集成 go-judge 进行真实评测。

## 核心功能概述
- 题库管理：题目创建、编辑、描述与样例管理。
- 比赛管理：比赛创建、题目编排、报名/退出、比赛期间排名与成绩展示。
- 提交评测：支持提交代码并通过 go-judge 执行编译与测试。
- 自测（Playground）：在线编写并自测代码，查看编译输出与运行结果。
- 管理后台：题目、比赛、公告、用户与提交管理界面（开发中）。
- 异步评测流水线：提交消息通过 RabbitMQ 下发给判题服务，评测结果异步写回并可查询。

## 技术栈概述
- 后端
  - golang (1.24) + Kratos（微服务框架）
  - ent ORM（维护实体/迁移）
  - gRPC / Protobuf（服务间 RPC）
  - RabbitMQ（消息队列）
  - Redis（缓存/短期状态）
  - MySQL（持久化）
  - go-judge（评测引擎）
- 前端
  - Vue3 + typescript + Vite
  - pnpm workspaces（monorepo 管理）
  - Monaco 编辑器（在线代码编辑）
  - marked + KaTeX（题目 Markdown 渲染与公式支持）
- DevOps / 部署
  - Docker（容器化与本地部署）
  - 项目提供部署脚本（支持 *NIX shell）
  - 各子服务拥有自己的构建 Dockerfile

## 项目结构
```bash
.
├── assets
├── cascade-oj
│   ├── api
│   ├── app
│   ├── cases
│   ├── config
│   ├── data
│   ├── ent
│   ├── pkg
│   ├── deploy.sh
│   ├── docker-compose.dependencies.yml
│   ├── docker-compose.yml
│   ├── Dockerfile
│   ├── generate_and_bundle_openapi.sh
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   ├── openapi.yaml
│   ├── openapi-merge.json
│   ├── README.md
│   └── table_creator.go
├── cascade-oj-frontend
│   ├── nginx
│   ├── packages
│   ├── public
│   ├── docker-compose.yml
│   ├── Dockerfile
│   ├── package.json
│   ├── package-lock.json
│   ├── pnpm-lock.yaml
│   ├── pnpm-workspace.yaml
│   ├── README.md
│   ├── setup.sh
│   └── tsconfig.json
├── LICENSE
└── README.md
```


## 第三方项目参考及使用

使用 go-judge 测评机镜像 [https://github.com/criyle/go-judge](https://github.com/criyle/go-judge)

微服务架构与项目结构参考了 SASTOJ [https://github.com/NJUPT-SAST/sastoj](https://github.com/NJUPT-SAST/sastoj)

## LICENSE

MIT license
