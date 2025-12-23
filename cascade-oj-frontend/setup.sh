#!/usr/bin/env bash
set -euo pipefail

# 用法:
#   ./setup.sh [options]
# 选项:
#   -h, --help       显示帮助
#   -f, --file FILE  使用自定义主 compose 文件（默认 ./docker-compose.yml）
#   --build          在 up 时构建镜像（等同 docker compose up --build）
#   --pull           先执行 pull 再 up
#   -n, --dry-run    只打印将要执行的命令，不实际运行
#
# 脚本行为:
# - 会优先使用系统的 `docker-compose`，若不可用则退回到 `docker compose`

echo "注意：运行 setup 脚本前需要手动构建前端资源"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_MAIN="$SCRIPT_DIR/docker-compose.yml"
BUILD=false
PULL=false
DRY_RUN=false

print_help() {
		cat <<'HELP'
用法:
  ./setup.sh [options]
选项:
  -h, --help       显示帮助
  -f, --file FILE  使用自定义主 compose 文件（默认 ./docker-compose.yml）
  --build          在 up 时构建镜像（等同 docker compose up --build）
  --pull           先执行 pull 再 up
  -n, --dry-run    只打印将要执行的命令，不实际运行

脚本行为:
- 会优先使用系统的 `docker-compose`，若不可用则退回到 `docker compose`

示例:
	# 部署默认 compose
	./setup.sh

	# 部署前端，并在启动前构建其镜像
	./setup.sh --build

	# 仅做 dry-run，查看将要执行的命令
	./setup.sh -n --build
HELP
}

detect_compose_cmd() {
	if command -v docker-compose >/dev/null 2>&1; then
		echo "docker-compose"
	elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
		# Use `docker compose` (note: it's a two-token command when invoked)
		echo "docker compose"
	else
		echo ""
	fi
}

# 解析参数
while [[ $# -gt 0 ]]; do
	case "$1" in
		-h|--help)
			print_help
			exit 0
			;;
		-f|--file)
			if [[ -z "${2-}" ]]; then
				echo "ERROR: -f|--file 需要一个参数" >&2
				exit 2
			fi
			COMPOSE_MAIN="$2"
			shift 2
			;;
		--build)
			BUILD=true
			shift
			;;
		--pull)
			PULL=true
			shift
			;;
		-n|--dry-run)
			DRY_RUN=true
			shift
			;;
		*)
			echo "未知参数: $1" >&2
			exit 2
			;;
	esac
done

DC_CMD="$(detect_compose_cmd)"
if [[ -z "$DC_CMD" ]]; then
	echo "ERROR: 未找到 'docker-compose' 或 'docker compose' 命令，请先安装 Docker CLI 或 docker-compose." >&2
	exit 3
fi

COMPOSE_ARGS=()
if [[ -f "$COMPOSE_MAIN" && -s "$COMPOSE_MAIN" ]]; then
	COMPOSE_ARGS+=("-f" "$COMPOSE_MAIN")
fi

if [[ ${#COMPOSE_ARGS[@]} -eq 0 ]]; then
	echo "ERROR: 未找到任何 compose 文件（默认查找 docker-compose.yml 与 docker-compose.dependencies.yml）。" >&2
	exit 4
fi

run_cmd() {
	echo "+ $*"
	if [[ "$DRY_RUN" == true ]]; then
		return 0
	fi
	# 如果 DC_CMD 包含空格（如 'docker compose'），需要用 eval 或数组展开
	if [[ "$DC_CMD" == "docker compose" ]]; then
		# join compose args safely
		docker compose "$@"
	else
		"$DC_CMD" "$@"
	fi
}

echo "使用的 compose 命令: $DC_CMD"
echo "compose 文件: ${COMPOSE_ARGS[*]}"

if [[ "$PULL" == true ]]; then
	echo "执行 docker pull..."
	if [[ "$DC_CMD" == "docker compose" ]]; then
		if [[ "$DRY_RUN" == true ]]; then
			echo "+ docker compose ${COMPOSE_ARGS[*]} pull"
		else
			docker compose "${COMPOSE_ARGS[@]}" pull
		fi
	else
		run_cmd "${COMPOSE_ARGS[@]}" pull
	fi
fi

# 构造 up 命令
UP_ARGS=("up" "-d")
if [[ "$BUILD" == true ]]; then
	UP_ARGS+=("--build")
fi

echo "执行: $DC_CMD ${COMPOSE_ARGS[*]} ${UP_ARGS[*]}"
if [[ "$DC_CMD" == "docker compose" ]]; then
	if [[ "$DRY_RUN" == true ]]; then
		echo "+ docker compose ${COMPOSE_ARGS[*]} ${UP_ARGS[*]}"
	else
		docker compose "${COMPOSE_ARGS[@]}" "${UP_ARGS[@]}"
	fi
else
	run_cmd "${COMPOSE_ARGS[@]}" "${UP_ARGS[@]}"
fi

echo "部署完成。可使用 'docker ps' 或 'docker compose ps' 查看运行状态。"
