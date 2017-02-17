#!/bin/sh

DB_HOST=${DB_HOST:-127.0.0.1}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-auth}
DB_USER=${DB_USER:-auth}
DB_PASS=${DB_PASS:-auth}

/kubernetes-authserver --host ${DB_HOST} --port ${DB_PORT} --db ${DB_NAME} --user ${DB_USER} --pass ${DB_PASS}