#!/usr/bin/env bash
# myplants 一键拉取 + 数据库备份(D1) + 构建 + 部署 (Gentoo Linux)
#
# 用法:
#   ./deploy.sh            # 默认 docker 模式: git pull + 复制 .env + 备份 DB 到 D1 + compose 构建滚动更新
#   ./deploy.sh native     # 裸机模式: 直接编译 Go + Vite 前端,OpenRC 重启
#   SKIP_D1=1 ./deploy.sh  # 跳过 D1 远程备份(仍保留本地备份)
#
# 可通过环境变量覆盖: APP_DIR / REPO_URL / BRANCH / DATA_DIR / DB_FILE / D1_DATABASE
# 流程: 拉代码 -> 把 $DATA_DIR/.env 复制到克隆目录 -> 备份 $DATA_DIR 下的 .db 到 D1 -> 构建部署
set -euo pipefail

APP_DIR="${APP_DIR:-/data/docker/myplants}"
REPO_URL="${REPO_URL:-https://github.zerlaer.cn//https://github.com/zerlaer/myplants.git}"
BRANCH="${BRANCH:-main}"
MODE="${1:-docker}"
HEALTH_URL="http://127.0.0.1:8020/"
DATA_DIR="${DATA_DIR:-/data/docker}"               # .env 与初始 myplants.db 所在目录
ENV_SRC="$DATA_DIR/.env"
DB_FILE="${DB_FILE:-}"                             # 留空则自动探测
D1_DATABASE="${D1_DATABASE:-myplants}"   # Cloudflare D1 库名
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
    log "工作区存在未提交修改:"
    git status --short
    if [ "${FORCE:-0}" = "1" ]; then
        log "FORCE=1: 自动 git stash 后继续(可用 git stash list 找回)"
        git stash push -m "deploy-auto-$(date +%Y%m%d-%H%M%S)"
    else
        err "请先处理上述文件后重试;确认服务器改动可丢弃/可暂存时用 FORCE=1 ./deploy.sh"
    fi
fi
git merge --ff-only "origin/$BRANCH"

# ---------- 2. 部署前置检查 ----------
# 从 DATA_DIR 复制 .env 到克隆目录(构建/运行所需凭据)
if [ -f "$ENV_SRC" ]; then
    log "复制 $ENV_SRC -> .env"
    cp -f "$ENV_SRC" .env
fi
[ -f .env ] || err "缺少 $ENV_SRC 及 $APP_DIR/.env(R2 凭证、CLOUDFLARE_API_TOKEN 等),参考 docker-compose.yaml 内注释创建"
set -a; . ./.env; set +a   # 让 wrangler 等子进程能读到 D1 所需变量

# ---------- 3. 部署前备份数据库(本地快照 + Cloudflare D1) ----------
# 自动探测 DB: 显式 DB_FILE > DATA_DIR/myplants.db > APP_DIR/data/myplants.db
if [ -z "$DB_FILE" ]; then
    if [ -f "$DATA_DIR/myplants.db" ]; then
        DB_FILE="$DATA_DIR/myplants.db"
    elif [ -f "$APP_DIR/data/myplants.db" ]; then
        DB_FILE="$APP_DIR/data/myplants.db"
    fi
fi
if [ -n "$DB_FILE" ] && [ -f "$DB_FILE" ]; then
    command -v sqlite3 >/dev/null 2>&1 || err "缺少 sqlite3 (emerge dev-db/sqlite),无法做一致性备份"
    TS=$(date +%Y%m%d-%H%M%S)
    BAK_DIR="$DATA_DIR/backups"
    mkdir -p "$BAK_DIR"
    BAK="$BAK_DIR/myplants-$TS.db"
    log "SQLite 一致性快照 $DB_FILE -> $BAK"
    sqlite3 "$DB_FILE" ".backup '$BAK'"

    if [ "${SKIP_D1:-0}" != "1" ]; then
        [ -n "${CLOUDFLARE_API_TOKEN:-}" ] || err "未设置 CLOUDFLARE_API_TOKEN,无法推送 D1 备份(可用 SKIP_D1=1 跳过)"
        [ -n "${CLOUDFLARE_ACCOUNT_ID:-}" ] || err "未设置 CLOUDFLARE_ACCOUNT_ID,wrangler 无法定位账户(在 .env 里补一行 CLOUDFLARE_ACCOUNT_ID=<账户ID>,或用 SKIP_D1=1 跳过)"
        DUMP="$BAK_DIR/myplants-$TS.sql"
        RAW="$BAK_DIR/myplants-$TS.raw.sql"
        # 不用 DROP 全库 schema(整库重建批次会被 D1 当作 reset 拒绝: D1_RESET_DO)
        # 改用: CREATE IF NOT EXISTS -> DELETE 清行 -> INSERT OR REPLACE
        sqlite3 "$BAK" ".dump" > "$RAW"
        {
            echo "PRAGMA foreign_keys=OFF;"
            sqlite3 "$BAK" "SELECT sql || ';' FROM sqlite_master WHERE type='table' AND sql IS NOT NULL AND name NOT LIKE 'sqlite_%';" \
                | sed -E 's/^CREATE TABLE /CREATE TABLE IF NOT EXISTS /I'
            sqlite3 "$BAK" "SELECT 'DELETE FROM ' || quote(name) || ';' FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';"
            { grep -E '^INSERT INTO ' "$RAW" || true; } | sed -E 's/^INSERT INTO /INSERT OR REPLACE INTO /'
        } > "$DUMP"
        rm -f "$RAW"
        log "推送到 Cloudflare D1: $D1_DATABASE ..."
        if npx --yes wrangler d1 execute "$D1_DATABASE" --remote --file "$DUMP" --yes; then
            rm -f "$DUMP"
            log "D1 备份完成 ✅"
        else
            err "D1 推送失败(本地快照与 $DUMP 已保留)。可 SKIP_D1=1 重跑,或检查 token/D1 库名"
        fi
    else
        log "已跳过 D1 远程备份"
    fi
    gzip -f "$BAK"
    # 本地只保留最近 $LOCAL_KEEP 份
    ls -t "$BAK_DIR"/myplants-*.db.gz 2>/dev/null | tail -n +$((LOCAL_KEEP + 1)) | xargs -r rm -f
else
    log "未找到数据库文件(DB_FILE / $DATA_DIR/myplants.db / $APP_DIR/data/myplants.db 均不存在,首次部署?),跳过备份"
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
