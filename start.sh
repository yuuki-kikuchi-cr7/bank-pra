
set -e

echo "run db migration"

# 環境変数の読み込みを確認
if [ -f /app/app.env ]; then
  echo "Loading environment variables from /app/app.env"
source /app/app.env
else
  echo "Error: /app/app.env not found"
  exit 1
fi

# DB_SOURCE が空の場合にエラーを出力
if [ -z "$DB_SOURCE" ]; then
  echo "Error: DB_SOURCE is not set"
  exit 1
fi
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"