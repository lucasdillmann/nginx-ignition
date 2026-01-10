#!/bin/sh
set -e

if command -v systemctl >/dev/null 2>&1; then
    if systemctl is-active --quiet nginx-ignition 2>/dev/null; then
        echo "Stopping nginx-ignition service..."
        systemctl stop nginx-ignition || true
    fi
    if systemctl is-enabled --quiet nginx-ignition 2>/dev/null; then
        echo "Disabling nginx-ignition service..."
        systemctl disable nginx-ignition || true
    fi
fi

exit 0
