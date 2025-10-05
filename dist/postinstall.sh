#!/bin/sh
set -e

if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
fi

echo ""
echo "nginx-ignition has been installed successfully!"
echo ""
echo "To start nginx-ignition, run:"
echo "  sudo systemctl enable nginx-ignition"
echo "  sudo systemctl start nginx-ignition"
echo ""
echo "Configuration file: /opt/nginx-ignition/nginx-ignition.properties"
echo "Service file: /etc/systemd/system/nginx-ignition.service"
echo ""
echo "After starting, access nginx-ignition at: http://localhost:8090"
echo ""

exit 0
