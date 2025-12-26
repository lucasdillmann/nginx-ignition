#!/bin/bash
set -e

echo "Starting backend development server..."
echo "Backend API will be available at http://localhost:8090"
echo ""
echo "Environment:"
echo "  - Database: postgres://localhost:5432/nginx_ignition"
echo "  - Database User: nginx_ignition"
echo "  - Database Password: devpassword"
echo ""

export NGINX_IGNITION_DATABASE_DRIVER=postgres
export NGINX_IGNITION_DATABASE_HOST=localhost
export NGINX_IGNITION_DATABASE_PORT=5432
export NGINX_IGNITION_DATABASE_NAME=nginx_ignition
export NGINX_IGNITION_DATABASE_SSL_MODE=disable
export NGINX_IGNITION_DATABASE_USERNAME=nginx_ignition
export NGINX_IGNITION_DATABASE_PASSWORD=devpassword
export NGINX_IGNITION_DATABASE_MIGRATIONS_PATH=/workspaces/nginx-ignition/database/common/migrations/scripts

cd /workspaces/nginx-ignition/application

if command -v air &> /dev/null; then
    echo "Using air for hot reload..."
    air
else
    echo "Running without hot reload"
    go run main.go
fi
