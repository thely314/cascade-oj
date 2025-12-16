#!/usr/bin/env bash
set -euo pipefail

# 可复用部署脚本
# 用法:
#   ./deploy.sh [options] [service1 service2 ...]
# 选项:
#   -h, --help       显示帮助
#   --no-deps        不包含 docker-compose.dependencies.yml
#   --deps           强制包含 docker-compose.dependencies.yml（默认行为）
#   -f, --file FILE  使用自定义主 compose 文件（默认 ./docker-compose.yml）
#   --build          在 up 时构建镜像（等同 docker compose up --build）
#   --pull           先执行 pull 再 up
#   -n, --dry-run    只打印将要执行的命令，不实际运行
#   --               后面的参数视为服务名
#
# 脚本行为:
# - 会优先使用系统的 `docker-compose`，若不可用则退回到 `docker compose`。
# - 如果存在 `docker-compose.dependencies.yml`，脚本默认会把它加入到 -f 列表中（可用 --no-deps 禁用）。

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COMPOSE_DEPS="$SCRIPT_DIR/docker-compose.dependencies.yml"
COMPOSE_MAIN="$SCRIPT_DIR/docker-compose.yml"
USE_DEPS=true
BUILD=false
PULL=false
DRY_RUN=false
DOWN=false
YES=false
REMOVE_VOLUMES=false
SERVICES=()

# 预设的有效服务列表（与 docker-compose.yml 中的 services 保持一致）
# 如果未来有更多服务，请同步更新此处或改为动态解析 compose 文件
PRESET_SERVICES=("gateway" "public" "user" "judge" "admin")

is_valid_service() {
	local svc="$1"
	for s in "${PRESET_SERVICES[@]}"; do
		if [[ "$s" == "$svc" ]]; then
			return 0
		fi
	done
	return 1
}

print_help() {
		cat <<'HELP'
可复用部署脚本
用法:
	./deploy.sh [options] [service1 service2 ...]

选项:
	-h, --help          显示帮助
	--no-deps           不包含 docker-compose.dependencies.yml
	--deps              强制包含 docker-compose.dependencies.yml（默认行为）
	-f, --file FILE     使用自定义主 compose 文件（默认 ./docker-compose.yml）
	--build             在 up 时构建镜像（等同 docker compose up --build）
	--pull              先执行 pull 再 up
	-n, --dry-run       只打印将要执行的命令，不实际运行
	--down              停止并移除容器（按服务或全部，不包括 dependencies），默认需要确认
	-y, --yes           非交互式确认（与 --down 一起使用）
	--volumes           与 --down 一起使用时同时移除卷（等同 -v）
	--                  后面的参数视为服务名

示例:
	# 部署默认 compose（包含 dependencies），后台运行全部服务
	./deploy.sh

	# 部署 gateway 服务，并在启动前构建其镜像
	./deploy.sh --build gateway

	# 仅做 dry-run，查看将要执行的命令
	./deploy.sh -n --build public

	# 非交互式下线全部服务并删除卷（CI/自动化）
	./deploy.sh --down -y --volumes

	# 交互式下线指定服务
	./deploy.sh --down public
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
		--no-deps)
			USE_DEPS=false
			shift
			;;
		--deps)
			USE_DEPS=true
			shift
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
		--down)
			DOWN=true
			shift
			;;
		-y|--yes)
			YES=true
			shift
			;;
		--volumes)
			REMOVE_VOLUMES=true
			shift
			;;
		--)
			shift
			while [[ $# -gt 0 ]]; do SERVICES+=("$1"); shift; done
			;;
		-*)
			echo "未知参数: $1" >&2
			exit 2
			;;
		*)
			SERVICES+=("$1")
			shift
			;;
	esac
done

DC_CMD="$(detect_compose_cmd)"
if [[ -z "$DC_CMD" ]]; then
	echo "ERROR: 未找到 'docker-compose' 或 'docker compose' 命令，请先安装 Docker CLI 或 docker-compose." >&2
	exit 3
fi

COMPOSE_ARGS=()
if [[ "$USE_DEPS" == true && -f "$COMPOSE_DEPS" ]]; then
	COMPOSE_ARGS+=("-f" "$COMPOSE_DEPS")
fi
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
if [[ ${#SERVICES[@]} -gt 0 ]]; then
	echo "要部署的服务: ${SERVICES[*]}"
else
	echo "将部署 compose 中定义的所有服务 (未指定服务名)"
fi

# 如果用户指定了服务名，校验它们是否在预设服务列表中
if [[ ${#SERVICES[@]} -gt 0 ]]; then
	INVALID=()
	for svc in "${SERVICES[@]}"; do
		if ! is_valid_service "$svc"; then
			INVALID+=("$svc")
		fi
	done

	if [[ ${#INVALID[@]} -gt 0 ]]; then
		echo "ERROR: 未知的服务: ${INVALID[*]}" >&2
		echo "可用的服务: ${PRESET_SERVICES[*]}" >&2
		exit 5
	fi
fi

	# 如果用户请求下线操作，执行安全的停止/移除流程
	if [[ "$DOWN" == true ]]; then
		if [[ "$DRY_RUN" == true ]]; then
			echo "+ DRY RUN: 将执行下线操作" 
		else
			# 交互确认，除非 YES=true
			if [[ "$YES" != true ]]; then
				echo "你即将停止并移除以下容器:"
				if [[ ${#SERVICES[@]} -gt 0 ]]; then
					echo "  ${SERVICES[*]}"
				else
					echo "  全部（使用 docker compose down）"
				fi
				read -r -p "确认继续？ (y/N): " ans
				case "$ans" in
					[Yy]|[Yy][Ee][Ss]) ;;
					*) echo "已取消"; exit 0 ;;
				esac
			fi

			if [[ ${#SERVICES[@]} -gt 0 ]]; then
				# 停止并移除指定服务
				echo "停止服务: ${SERVICES[*]}"
				run_cmd "${COMPOSE_ARGS[@]}" stop "${SERVICES[@]}"
				RM_ARGS=("rm" "-s" "-f")
				if [[ "$REMOVE_VOLUMES" == true ]]; then
					RM_ARGS+=("-v")
				fi
				run_cmd "${COMPOSE_ARGS[@]}" "${RM_ARGS[@]}" "${SERVICES[@]}"
			else
				# 对整个 compose 做 down
				DOWN_ARGS=("down")
				if [[ "$REMOVE_VOLUMES" == true ]]; then
					DOWN_ARGS+=("-v")
				fi
				run_cmd "${COMPOSE_ARGS[@]}" "${DOWN_ARGS[@]}"
			fi
		fi
		echo "下线完成。"
		exit 0
	fi

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

if [[ ${#SERVICES[@]} -gt 0 ]]; then
	UP_ARGS+=("${SERVICES[@]}")
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
