#!/bin/bash

# .env ファイルを読み込む
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# ENVが設定されているか確認する
if [ -z "$ENV" ]; then
  echo "Error: ENV is not set. Please specify an environment (development, staging, production)."
  exit 1
fi

# DB_HOST にデフォルト値を指定（デフォルトは localhost）
POSTGRES_HOST=${POSTGRES_HOST:-localhost}

# ENVに対応するディレクトリのパスを設定
SEED_DIR="$(dirname "$0")/$ENV"

# seeds ディレクトリにあるすべての .sql ファイルを対象にする
SEED_FILES=$(ls $SEED_DIR/*.sql)

# 各 SQL ファイルをデータベースに適用する
for file in $SEED_FILES; do
  echo "Applying $file..."
  psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -w -f $file --single-transaction
done
