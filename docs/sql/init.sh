#!/bin/bash

set -euxo pipefail

HOST="${MYSQL_HOST:-127.0.0.1}"
PORT="${MYSQL_PORT:-3306}"
DB_NAME="${MYSQL_DB:-gg_project}"
USERNAME="${MYSQL_USER:-root}"
PASSWORD="${MYSQL_PASSWORD:-password}"
PROTOCOL="${MYSQL_PROTOCOL:-tcp}"

rm temp.sql || true
echo "DROP DATABASE IF EXISTS ${DB_NAME};\n" >> temp.sql
echo "CREATE DATABASE IF NOT EXISTS ${DB_NAME};\n" >> temp.sql
echo "USE ${DB_NAME};\n" >> temp.sql
echo "SET foreign_key_checks = 0;\n" >> temp.sql
cat ./[!temp]*.sql >> temp.sql
echo "SET foreign_key_checks = 1;\n"

mysql -h ${HOST} -P ${PORT} -u ${USERNAME} --password=${PASSWORD} --protocol=${PROTOCOL} < temp.sql
rm temp.sql
