#!/bin/zsh
# NoSocial 服务管理脚本
# 用法: ./nos.sh [start|stop|status|restart]

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
FRONTEND_DIR="$SCRIPT_DIR/frontend"
LOG_DIR="$SCRIPT_DIR/backend/logs"

# ============================================================
# 后端配置环境变量（NOSOCIAL_ 前缀，由 viper AutomaticEnv 映射）
# 优先级：外部已存在的环境变量  >  backend/.env 文件  >  脚本内兜底默认值
# config.yaml 中对应字段应留空，避免明文泄露
# ============================================================

ENV_FILE="$BACKEND_DIR/.env"
if [ -f "$ENV_FILE" ]; then
  # 逐行读取 .env：跳过空行与注释行，支持值内含 = 号；
  # 仅当该变量尚未在外部环境中定义时才导出，保证"外部 > .env > 默认值"的优先级。
  while IFS= read -r line || [ -n "$line" ]; do
    # 去掉行首尾空白
    line="${line#"${line%%[![:space:]]*}"}"
    line="${line%"${line##*[![:space:]]}"}"
    [ -z "$line" ] && continue
    case "$line" in
      \#*) continue ;;
    esac
    key="${line%%=*}"
    value="${line#*=}"
    # 去除 value 两端的成对引号
    case "$value" in
      \"*\") value="${value#\"}"; value="${value%\"}" ;;
      \'*\') value="${value#\'}"; value="${value%\'}" ;;
    esac
    # 仅在环境中尚未定义时注入
    if [ -z "${(P)key-}" ]; then
      export "$key=$value"
    fi
  done < "$ENV_FILE"
fi

# --- App ---
export NOSOCIAL_APP_MODE="${NOSOCIAL_APP_MODE:-debug}"
export NOSOCIAL_APP_PORT="${NOSOCIAL_APP_PORT:-8080}"

# --- PostgreSQL ---
export NOSOCIAL_POSTGRESQL_HOST="${NOSOCIAL_POSTGRESQL_HOST:-123.57.241.163}"
export NOSOCIAL_POSTGRESQL_PORT="${NOSOCIAL_POSTGRESQL_PORT:-5432}"
export NOSOCIAL_POSTGRESQL_USER="${NOSOCIAL_POSTGRESQL_USER:-postgres}"
export NOSOCIAL_POSTGRESQL_PASSWORD="${NOSOCIAL_POSTGRESQL_PASSWORD:-Nzi2001oak*.}"
export NOSOCIAL_POSTGRESQL_DBNAME="${NOSOCIAL_POSTGRESQL_DBNAME:-nosocial}"

# --- Redis ---
export NOSOCIAL_REDIS_HOST="${NOSOCIAL_REDIS_HOST:-123.57.241.163}"
export NOSOCIAL_REDIS_PORT="${NOSOCIAL_REDIS_PORT:-6379}"
export NOSOCIAL_REDIS_PASSWORD="${NOSOCIAL_REDIS_PASSWORD:-000000}"
export NOSOCIAL_REDIS_DB="${NOSOCIAL_REDIS_DB:-0}"

# --- JWT Secret ---
export NOSOCIAL_JWT_ADMIN_SECRET="${NOSOCIAL_JWT_ADMIN_SECRET:-lK4OgH9vnwRefKKDJ5tPkaw7vSoZTLrenk7FoSVlRU0=}"
export NOSOCIAL_JWT_WX_SECRET="${NOSOCIAL_JWT_WX_SECRET:-HT7nIkRlymvngokQDdzuueSXHqZ5+5ZG4hRi/TjkV5U=}"

# --- OSS ---
export NOSOCIAL_OSS_ACCESS_KEY_ID="${NOSOCIAL_OSS_ACCESS_KEY_ID:-LTAI5t89vpzWokSwxJJTeCzT}"
export NOSOCIAL_OSS_ACCESS_KEY_SECRET="${NOSOCIAL_OSS_ACCESS_KEY_SECRET:-93MzcnCso7QXfZKVj3SybYjJZaPgqw}"

# --- WX（小程序登录/支付，值统一从 backend/.env 读取；未配置则留空） ---
export NOSOCIAL_WX_APPID="${NOSOCIAL_WX_APPID:-}"
export NOSOCIAL_WX_SECRET="${NOSOCIAL_WX_SECRET:-}"
export NOSOCIAL_WX_API_V3_KEY="${NOSOCIAL_WX_API_V3_KEY:-}"

# 颜色
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'

log()   { echo -e "${GREEN}[NoSocial]${NC} $1"; }
info()  { echo -e "${CYAN}[INFO]${NC}  $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $1"; }
die()   { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# 获取端口对应的 PID
get_pid_by_port() {
  lsof -i :$1 -sTCP:LISTEN -t 2>/dev/null | head -1
}

# 检查端口是否被占用
is_port_used() {
  [ -n "$(get_pid_by_port $1)" ]
}

# 等待端口就绪
wait_for_port() {
  local port=$1; local name=$2
  for i in {1..15}; do
    curl -s http://localhost:$port/api/v1/public/shop -o /dev/null 2>/dev/null && return 0
    sleep 1
  done
  return 1
}

# ---------- start ----------
do_start() {
  log "启动 NoSocial 服务..."

  # 停止已有进程
  killall nosocial 2>/dev/null || true
  killall node 2>/dev/null || true
  sleep 1

BIN_DIR="$SCRIPT_DIR/bin"
mkdir -p "$BIN_DIR"

  # 编译后端（每次启动都重新编译以确保带上最新代码）
  cd "$BACKEND_DIR"
  info "编译后端 → bin/nosocial"
  go build -o "$BIN_DIR/nosocial" ./cmd/server/ || die "后端编译失败"

  # 启动后端
  mkdir -p "$LOG_DIR"
  "$BIN_DIR/nosocial" > "$LOG_DIR/nosocial.log" 2>&1 &
  local backend_pid=$!
  log "后端已启动 (PID: $backend_pid)"

  # 等待后端就绪
  sleep 3
  if ! wait_for_port 8080 "后端"; then
    die "后端启动超时，请查看 logs/nosocial.log"
  fi
  log "后端就绪 ✓"

  # 启动前端
  cd "$FRONTEND_DIR"
  npx vite --host 0.0.0.0 --port 5173 > /dev/null 2>&1 &
  local frontend_pid=$!
  log "前端已启动 (PID: $frontend_pid)"

  # 等待前端就绪
  sleep 5
  if ! wait_for_port 5173 "前端"; then
    warn "前端可能未就绪，请稍后重试"
  else
    log "前端就绪 ✓"
  fi

  echo ""
  echo -e "${GREEN}======================================${NC}"
  echo -e "${GREEN}  NoSocial 服务已全部启动${NC}"
  echo -e "  后端 API: ${YELLOW}http://localhost:8080/api/v1/${NC}"
  echo -e "  管理端:   ${YELLOW}http://localhost:5173/${NC}"
  echo -e "${GREEN}======================================${NC}"
  echo "  后端 PID: $backend_pid"
  echo "  前端 PID: $frontend_pid"
  echo ""
}

# ---------- stop ----------
do_stop() {
  log "停止 NoSocial 服务..."
  local count=0

  for pid in $(lsof -i :8080 -t 2>/dev/null); do
    kill -9 $pid 2>/dev/null && info "已终止后端进程 $pid" || true
    ((count++)) || true
  done

  for pid in $(lsof -i :5173 -t 2>/dev/null); do
    kill -9 $pid 2>/dev/null && info "已终止前端进程 $pid" || true
    ((count++)) || true
  done

  killall nosocial 2>/dev/null || true
  killall -9 nosocial 2>/dev/null || true

  if [ $count -gt 0 ] || pgrep -x nosocial >/dev/null 2>&1; then
    log "所有服务已停止 ✓"
  else
    info "没有运行中的服务"
  fi
}

# ---------- status ----------
do_status() {
  echo -e "${CYAN}=== NoSocial 服务状态 ===${NC}"
  echo ""

  local backend_ok=false; local frontend_ok=false

  if is_port_used 8080; then
    local bpid=$(get_pid_by_port 8080)
    echo -e "  ${GREEN}●${NC} 后端 Go API  : ${GREEN}运行中${NC}  (PID: $bpid, 端口: 8080)"
    backend_ok=true
  else
    echo -e "  ${RED}○${NC} 后端 Go API  : ${RED}未运行${NC}  (端口: 8080)"
  fi

  if is_port_used 5173; then
    local fpid=$(get_pid_by_port 5173)
    echo -e "  ${GREEN}●${NC} 前端 React    : ${GREEN}运行中${NC}  (PID: $fpid, 端口: 5173)"
    frontend_ok=true
  else
    echo -e "  ${RED}○${NC} 前端 React    : ${RED}未运行${NC}  (端口: 5173)"
  fi

  echo ""
  if $backend_ok && $frontend_ok; then
    echo -e "  状态: ${GREEN}全部在线${NC}"
    echo -e "  后端: ${YELLOW}http://localhost:8080/api/v1/${NC}"
    echo -e "  前端: ${YELLOW}http://localhost:5173/${NC}"
  elif ! $backend_ok && ! $frontend_ok; then
    echo -e "  状态: ${RED}全部离线${NC}"
  else
    echo -e "  状态: ${YELLOW}部分在线${NC}"
  fi
  echo ""
}

# ---------- restart ----------
do_restart() {
  do_stop
  sleep 2
  do_start
}

# ---------- main ----------
CMD="${1:-start}"
case "$CMD" in
  start|stop|status|restart)
    do_$CMD
    ;;
  *)
    echo "用法: $0 {start|stop|status|restart}"
    exit 1
    ;;
esac
