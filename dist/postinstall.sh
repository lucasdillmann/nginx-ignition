#!/bin/sh
set -e

if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || true
fi

echo ""
echo "nginx-ignition has been installed successfully"
echo ""
echo "To start nginx-ignition, run:"
echo "  sudo systemctl start nginx-ignition"
echo ""
echo "If you wish for nginx-ignition to start on boot, run:"
echo "  sudo systemctl enable nginx-ignition"
echo ""
echo "After starting, nginx-ignition will be available at http://localhost:8090"
echo "You can change the port by editing the configuration file at /opt/nginx-ignition/nginx-ignition.properties"
echo ""
echo "Please note that nginx-ignition starts by default using an embedded SQLite database. "
echo "For production environments, it is recommended to use a PostgreSQL instead. Please refer to the"
echo "documentation at https://github.com/lucasdillmann/nginx-ignition for more information on how to make the change."
echo ""

exit 0
