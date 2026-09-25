#!/usr/bin/env bash
# 把服务器上当前的 SQLite 数据库备份到 Cloudflare D1(独立脚本,可手动/定时执行)
#
# 用法:
#   ./backup-to-d1.sh                      # 自动探测数据库文件并备份
#   DB_FILE=/data/docker/myplants.db ./backup-to-d1.sh   # 显式指定
#
# 依赖: sqlite3、node/npx(wrangler)、.env 里的 CLOUDFLARE_API_TOKEN
set -euo pipefail

DATA_DIR="${DATA_DIR:-/data/docker}"
APP_DIR="${APP_DIR:-$DATA_DIR/myplants}"
DB_FILE="${DB_FILE:-}"
D1_DATABASE="${D1_DATABASE:-myplants}"
BAK_DIR="${BAK_DIR:-$DATA_DIR/backups}"
LOCAL_KEEP=7

log() { printf '\e[32m[backup]\e[0m %s\n' "$*"; }
err() { printf '\e[31m[backup]\e[0m %s\n' "$*" >&2; exit 1; }

# ---------- 1. 凭据 ----------
for f in "$APP_DIR/.env" "$DATA_DIR/.env"; do
    if [ -f "$f" ]; then
        set -a; . "$f"; set +a
        log "已加载 $f"
        break
    fi
done
[ -n "${CLOUDFLARE_API_TOKEN:-}" ] || err "未找到 CLOUDFLARE_API_TOKEN(请在 .env 中配置)"
[ -n "${CLOUDFLARE_ACCOUNT_ID:-}" ] || err "未找到 CLOUDFLARE_ACCOUNT_ID: 在 .env 里加一行 CLOUDFLARE_ACCOUNT_ID=<账户ID> (Cloudflare 控制台右侧账户ID,与 R2 账号同一个)"

# ---------- 2. 探测数据库文件 ----------
if [ -z "$DB_FILE" ]; then
    if [ -f "$DATA_DIR/myplants.db" ]; then
        DB_FILE="$DATA_DIR/myplants.db"
    elif [ -f "$APP_DIR/data/myplants.db" ]; then
        DB_FILE="$APP_DIR/data/myplants.db"
    else
        first=$(ls "$DATA_DIR"/*.db 2>/dev/null | head -1 || true)
        [ -n "$first" ] && DB_FILE="$first"
    fi
fi
[ -n "$DB_FILE" ] && [ -f "$DB_FILE" ] || err "找不到数据库文件,可用 DB_FILE=/path/to.db 指定"
log "源数据库: $DB_FILE ($(du -h "$DB_FILE" | cut -f1))"

command -v sqlite3 >/dev/null 2>&1 || err "缺少 sqlite3 (emerge dev-db/sqlite)"

# ---------- 3. 一致性快照(避免拷到写了一半的库) ----------
TS=$(date +%Y%m%d-%H%M%S)
mkdir -p "$BAK_DIR"
BAK="$BAK_DIR/myplants-$TS.db"
sqlite3 "$DB_FILE" ".backup '$BAK'"
log "快照 -> $BAK"

# ---------- 4. 生成 SQL 并推送 D1 ----------
# 不用 DROP 全库 schema(整库重建批次会被 D1 当作 reset 拒绝: D1_RESET_DO)
# 改用: CREATE IF NOT EXISTS -> DELETE 清行 -> INSERT OR REPLACE
DUMP="$BAK_DIR/myplants-$TS.sql"
RAW="$BAK_DIR/myplants-$TS.raw.sql"
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
    err "D1 推送失败(本地快照与 $DUMP 已保留,可修正后手动重跑 wrangler),请检查 token / D1 库名"
fi

# ---------- 5. 本地保留最近 N 份 ----------
gzip -f "$BAK"
ls -t "$BAK_DIR"/myplants-*.db.gz 2>/dev/null | tail -n +$((LOCAL_KEEP + 1)) | xargs -r rm -f
log "本地快照: $BAK.gz (保留最近 $LOCAL_KEEP 份)"
