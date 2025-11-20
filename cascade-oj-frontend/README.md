# Cascade-oj-frontend

Cascade-oj 的前端项目

前端采用 pnpm workspaces 组织为 monorepo 结构

在 cascade-oj-frontend/packages 目录下存放子项目

## Setup environment

```bash
# At cascade-oj/cascade-oj-frontend
pnpm install
```

## Dev

```bash
# At cascade-oj/cascade-oj-frontend
pnpm run project_name:dev
# example
pnpm run competition:dev
```

## Build

```bash
# At cascade-oj/cascade-oj-frontend
pnpm run project_name:build
# example
pnpm run competition:build
```

如果你想直接在子包内执行命令，也可以使用以下命令：
```bash
pnpm --filter <package name> <command>
```
