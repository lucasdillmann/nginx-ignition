#!/bin/sh

mkdir -p /var/lib/tailscale
chmod 700 /var/lib/tailscale
/usr/sbin/tailscaled \
  --state=mem:nginx-ignition \
  --statedir=/tmp/nginx-ignition/tailscale \
  --socket=/var/lib/tailscale/tailscaled.sock \
  --verbose -1 >/dev/null 2>&1 &

exec /opt/nginx-ignition/nginx-ignition
