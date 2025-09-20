#!/bin/bash
set -e

DB_NAME="${POSTGRES_DB}_test"

# Check if database exists; create if not (Postgres has no CREATE DATABASE IF NOT EXISTS)
if psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -tAc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}';" | grep -q 1; then
    echo "Database ${DB_NAME} already exists"
else
    psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c "CREATE DATABASE \"${DB_NAME}\";"
fi