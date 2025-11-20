# Cascade-oj

Cascade-oj 的后端项目，一个微服务 online judge 系统，基于 kratos

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
