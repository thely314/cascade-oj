# Cascade-oj

Cascade-oj 的后端项目，一个微服务 online judge 系统，基于 kratos

data/ 目录用于挂载容器持久化，如数据库、日志落盘
cases/ 目录用于存放题目测试用例，同样挂载到容器，具体结构查阅 cases/README.md

## Dev
### Environment

```
go version: 1.24
kratos version: v2.9
```

### Dependency & Build

后端部分使用 Makefile 控制构建，请确保拥有 make
对子服务的递归构建（如 config）直接在 shell 调用了 make，因此 mingw32-make 用户可能需要使用别名

具体操作可使用 `make help` 查看，使用 make 时注意工作目录或显式指定的 Makefile

> 经过测试，很难在 Windows powershell 环境下获取 git-bash 路径并完成传参，所以还是要求开发人员**直接使用 git-bash 等 unix shell 执行 make 操作**

```bash
# At cascade-oj/cascade-oj
make help
```

更详细的指引请看 Makefile

### Docker
推荐使用 Docker 部署，项目提供了相应的 Dockerfile 和 docker-compose

## Deploy

仓库根目录下包含一个可复用的部署脚本 `deploy.sh`，以及示例 `docker-compose.yml` 和 `docker-compose.dependencies.yml`（后者包含 MySQL 等依赖）。脚本支持按服务部署或全部部署，并可选择在部署时构建镜像。

主要功能与参数：

- `-h, --help`：显示帮助
- `--no-deps`：不包含 `docker-compose.dependencies.yml`（默认会包含）
- `-f, --file FILE`：使用自定义 compose 文件（默认 `./docker-compose.yml`）
- `--build`：在 up 时构建镜像（等同 `docker compose up --build`），会使用各子服务目录下的 `Dockerfile`（如果存在）作为 builder
- `--pull`：先执行 pull 再启动
- `-n, --dry-run`：只打印将要执行的命令，不实际运行
- `--down`：停止并移除容器（可指定服务或全部，不包括 dependencies）。默认会进行交互式确认以避免误操作。
- `-y, --yes`：与 `--down` 一起使用，跳过交互确认（适合 CI/自动化脚本）。
- `--volumes`：与 `--down` 一起使用时同时移除卷（等同 `-v`）。

示例：

在项目根目录下（即包含 `deploy.sh` 的目录）执行：

```bash
# 部署默认 compose（包含 dependencies），后台运行全部服务
./deploy.sh

# 部署 gateway 服务，并在启动前构建其镜像
./deploy.sh --build gateway

# 仅做 dry-run，查看将要执行的命令
./deploy.sh -n --build public

# 指定自定义 compose 文件，不包含依赖 compose
./deploy.sh -f docker-compose.prod.yml --no-deps --pull

# 交互式下线全部服务
./deploy.sh --down

# 非交互式下线并删除卷（可在 CI 中使用）
./deploy.sh --down -y --volumes

# 下线单个服务
./deploy.sh --down public
```

关于 builder:

每个微服务目录（例如 `app/gateway`、`app/services/public`）都可以包含自己的 `Dockerfile`，作为该服务的构建器。`deploy.sh --build` 会触发 compose 在这些服务目录中执行镜像构建（等同 `docker compose build` 或 `docker compose up --build`），所以运维可以选择在部署时构建最新的服务镜像，或先 `--pull` 使用远程镜像。

建议：在生产环境中，推荐在 CI 中构建并推送镜像（避免在部署时编译），在运维侧使用 `./deploy.sh --pull` 来拉取已发布的镜像并启动。
