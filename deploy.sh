#!/usr/bin/env bash
# myplants 一键拉取 + 数据库备份(D1) + 构建 + 部署 (Gentoo Linux)
#
# 用法:
#   ./deploy.sh            # 默认 docker 模式: git pull + 备份 DB + compose 构建滚动更新
#   ./deploy.sh native     # 裸机模式: 直接编译 Go + Vite 前端,OpenRC 重启
#   SKIP_D1=1 ./deploy.sh  # 跳过 D1 远程备份(仍保留本地备份)
#
# 可通过环境变量覆盖: APP_DIR / REPO_URL / BRANCH / D1_DATABASE
set -euo pipefail

APP_DIR="${APP_DIR:-/data/docker/myplants}"
REPO_URL="${REPO_URL:-https://github.com/zerlaer/myplants.git}"
BRANCH="${BRANCH:-main}"
MODE="${1:-docker}"
HEALTH_URL="http://127.0.0.1:8020/"
DB_FILE="$APP_DIR/data/myplants.db"
D1_DATABASE="${D1_DATABASE:-myplants-backup}"   # Cloudflare D1 库名
LOCAL_KEEP=7                                    # 本地备份保留份数

log() { printf '\e[32m[deploy]\e[0m %s\n' "$*"; }
err() { printf '\e[31m[deploy]\e[0m %s\n' "$*" >&2; exit 1; }

# ---------- 1. 拉取代码 ----------
if [ ! -d "$APP_DIR/.git" ]; then
    log "首次克隆 $REPO_URL -> $APP_DIR"
    git clone -b "$BRANCH" "$REPO_URL" "$APP_DIR"
fi
cd "$APP_DIR"
log "更新分支 $BRANCH ..."
git fetch origin "$BRANCH"
if ! git diff --quiet || ! git diff --cached --quiet; then
    err "工作区存在未提交修改,拒绝覆盖。请先 git stash/commit 或手动部署"
fi
git merge --ff-only "origin/$BRANCH"

# ---------- 2. 部署前置检查 ----------
[ -f .env ] || err "缺少 $APP_DIR/.env(R2 凭证、CLOUDFLARE_API_TOKEN 等),参考 docker-compose.yaml 内注释创建"
set -a; . ./.env; set +a   # 让 wrangler 等子进程能读到 D1 所需变量

# ---------- 3. 部署前备份数据库(本地快照 + Cloudflare D1) ----------
if [ -f "$DB_FILE" ]; then
    command -v sqlite3 >/dev/null 2>&1 || err "缺少 sqlite3 (emerge dev-db/sqlite),无法做一致性备份"
    TS=$(date +%Y%m%d-%H%M%S)
    BAK_DIR="$APP_DIR/data/backups"
    mkdir -p "$BAK_DIR"
    BAK="$BAK_DIR/myplants-$TS.db"
    log "SQLite 一致性快照 -> $BAK"
    sqlite3 "$DB_FILE" ".backup '$BAK'"

    if [ "${SKIP_D1:-0}" != "1" ]; then
        [ -n "${CLOUDFLARE_API_TOKEN:-}" ] || err "未设置 CLOUDFLARE_API_TOKEN,无法推送 D1 备份(可用 SKIP_D1=1 跳过)"
        DUMP="$BAK_DIR/myplants-$TS.sql"
        {
            echo "PRAGMA foreign_keys=OFF;"
            sqlite3 "$BAK" "SELECT 'DROP TABLE IF EXISTS ' || quote(name) || ';' FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';"
            sqlite3 "$BAK" ".dump"
        } > "$DUMP"
        log "推送到 Cloudflare D1: $D1_DATABASE ..."
        if npx --yes wrangler d1 execute "$D1_DATABASE" --remote --file "$DUMP" --yes; then
            rm -f "$DUMP"
            log "D1 备份完成 ✅"
        else
            err "D1 推送失败(本地快照已保留)。可 SKIP_D1=1 重跑,或检查 token/D1 库名"
        fi
    else
        log "已跳过 D1 远程备份"
    fi
    gzip -f "$BAK"
    # 本地只保留最近 $LOCAL_KEEP 份
    ls -t "$BAK_DIR"/myplants-*.db.gz 2>/dev/null | tail -n +$((LOCAL_KEEP + 1)) | xargs -r rm -f
else
    log "未找到 $DB_FILE(首次部署?),跳过备份"
fi

case "$MODE" in

# ---------- 4a. Docker 模式 ----------
docker)
    command -v docker >/dev/null 2>&1 || err "未安装 docker (app-emulation/docker 或 docker-ce)"
    docker info >/dev/null 2>&1 || { log "启动 docker 服务"; rc-service docker start; }
    COMPOSE="docker compose"
    docker compose version >/dev/null 2>&1 || COMPOSE="docker-compose"
    log "构建镜像(拉取基础镜像更新)..."
    $COMPOSE build --pull
    log "滚动更新容器..."
    $COMPOSE up -d
    $COMPOSE image prune -f >/dev/null 2>&1 || true
    ;;

# ---------- 4b. 裸机模式 ----------
native)
    command -v go >/dev/null 2>&1 || err "未安装 dev-lang/go (需 >=1.21)"
    command -v npm >/dev/null 2>&1 || err "未安装 net-libs/nodejs"
    log "构建前端..."
    (cd frontend && npm ci && npm run build)
    log "构建后端..."
    go build -trimpath -o myplants .
    if [ -f /etc/init.d/myplants ]; then
        log "重启 OpenRC 服务..."
        rc-service myplants restart
    else
        err "已编译完成,但缺少 /etc/init.d/myplants 服务脚本(模板见本文件末尾注释)"
    fi
    ;;

*)
    err "用法: $0 [docker|native]"
    ;;
esac

# ---------- 5. 健康检查 ----------
log "等待服务就绪..."
for _ in $(seq 1 30); do
    if wget -q -O /dev/null --timeout=3 "$HEALTH_URL" 2>/dev/null; then
        log "部署完成,服务已响应 ✅"
        exit 0
    fi
    sleep 2
done
err "健康检查超时。排查: docker logs --tail 50 myplants 或 tail logs/myplants.log"

# ---------- native 模式 OpenRC 服务脚本模板 ----------
# 保存到 /etc/init.d/myplants 并 rc-update add myplants default:
#
#   #!/sbin/openrc-run
#   name="myplants"
#   command="/data/docker/myplants/myplants"
#   command_background="yes"
#   pidfile="/run/${SVCNAME}.pid"
#   supervisor="supervise-daemon"
#   output_log="/var/log/myplants.log"
#   error_log="/var/log/myplants.err.log"
#   depend() { need net; after docker; }
#
# 然后: chmod +x /etc/init.d/myplants
