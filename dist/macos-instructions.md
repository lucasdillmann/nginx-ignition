# macOS installation and upgrade instructions

These instructions explain how to install and run nginx-ignition on a macOS host using launchd. The distribution .zip
you downloaded already includes:

- nginx-ignition binary
- frontend assets (UI)
- database migrations
- default configuration file, nginx-ignition.properties
- a ready-to-use launchd plist, nginx-ignition.plist

## Requirements
- macOS 12 or newer (Apple Silicon arm64)
- nginx installed and accessible (default expected path is /usr/sbin/nginx; can be changed in the properties file)
  - Common Homebrew paths: /opt/homebrew/bin/nginx (Apple Silicon)
- nginx modules mod-http-js, mod-http-lua and mod-stream installed
- Administrator (sudo) privileges for installation
- Port 8090 available (default UI/API port; also configurable)

## New installations
### Unpack the distribution
Unzip the archive you built or downloaded and copy its contents to /opt/nginx-ignition.

```bash
unzip nginx-ignition.macos-<arch>.zip -d /tmp/nginx-ignition-dist
sudo mkdir -p /opt/nginx-ignition
sudo cp -r /tmp/nginx-ignition-dist/* /opt/nginx-ignition/
```

Optional: customize and install the provided launchd plist example file
```bash
sudo cp /tmp/nginx-ignition-dist/nginx-ignition.plist /Library/LaunchDaemons/
```

### Fix the permissions
```bash
sudo chown -R root:wheel /opt/nginx-ignition
sudo chmod 0755 /opt/nginx-ignition/nginx-ignition
sudo chmod 0644 /opt/nginx-ignition/nginx-ignition.properties
sudo chown root:wheel /Library/LaunchDaemons/nginx-ignition.plist 2>/dev/null || true
sudo chmod 0644 /Library/LaunchDaemons/nginx-ignition.plist 2>/dev/null || true
```

If you prefer to run as a non-root user, remember to also change the UserName in /Library/LaunchDaemons/nginx-ignition.plist,
set appropriate ownership of /opt/nginx-ignition and /tmp/nginx-ignition, and grant the needed permissions for the user
to interact with the nginx binary and alike.

### Configure nginx ignition
Open /opt/nginx-ignition/nginx-ignition.properties and adjust values as needed. Please check the documentation
available on GitHub for full details, default values and recommendations.

Notable setting for macOS/Homebrew nginx:
- nginx-ignition.nginx.binary-path=/opt/homebrew/bin/nginx

nginx ignition can also be configured using only environment variables, removing the need for the configuration 
properties file entirely. The keys are the same as in the properties file, but with uppercase names and
dots/dashes replaced by underscores. For example, the `nginx-ignition.nginx.binary-path` property can be set 
using the `NGINX_IGNITION_NGINX_BINARY_PATH` environment variable.

### Register and start the service (using launchd)
Load the daemon so it starts at boot and immediately:

```bash
sudo launchctl load -w /Library/LaunchDaemons/nginx-ignition.plist
sudo launchctl start nginx-ignition
```

### Check status and logs
```bash
launchctl list | grep nginx-ignition
sudo launchctl print system/nginx-ignition
log stream --predicate 'process == "nginx-ignition"' --style syslog
```

### Open firewall (if applicable)
- If using the macOS Application Firewall, allow incoming connections for nginx-ignition via System Settings > 
  Network > Firewall > Options.
- If you use a third-party firewall, allow TCP port 8090.

### Open nginx ignition in your browser and start using
Open http://<your-host>:8090 in your browser. The nginx ignition UI will guide you through the remaining setup steps.

## Upgrading
1. Stop service: `sudo launchctl stop nginx-ignition && sudo launchctl unload -w /Library/LaunchDaemons/nginx-ignition.plist`
2. Replace /opt/nginx-ignition/nginx-ignition and any updated assets (frontend, migrations)
3. Start service: `sudo launchctl load -w /Library/LaunchDaemons/nginx-ignition.plist && sudo launchctl start nginx-ignition`
