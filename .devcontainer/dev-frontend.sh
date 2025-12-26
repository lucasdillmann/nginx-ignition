#!/bin/bash
set -e

echo "Starting frontend development server..."
echo "Frontend will be available at http://localhost:5173"
echo "API proxy target: http://localhost:8090"
echo ""

# Navigate to frontend directory
cd /workspaces/nginx-ignition/frontend

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

# Start Vite dev server
npm run start
