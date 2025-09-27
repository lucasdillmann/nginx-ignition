# Linux installation and upgrade instructions

These instructions explain how to install and run nginx-ignition on a Linux host using systemd. The distribution .zip 
you downloaded already includes:

- nginx-ignition binary
- frontend assets (UI)
- database migrations
- default configuration file, nginx-ignition.properties
- a ready-to-use systemd unit, nginx-ignition.service

## Requirements
- Linux (amd64/x86_64 or arm64)
- nginx installed and accessible (default expected path is /usr/sbin/nginx; can be changed in the properties file)
- nginx modules mod-http-js, mod-http-lua and mod-stream installed
- Root/sudo privileges for installation
- Port 8090 available (default UI/API port; also configurable)

## New installations
### Unpack the distribution
Unzip the archive you built or downloaded (linux-amd64 or linux-arm64) and copy its contents to /opt/nginx-ignition.

```bash
unzip nginx-ignition.linux-<arch>.zip -d /tmp/nginx-ignition-dist
sudo mkdir -p /opt/nginx-ignition
sudo cp -r /tmp/nginx-ignition-dist/* /opt/nginx-ignition/
```

Optional: customize and install the provided systemd unit example file
```bash
sudo cp /tmp/nginx-ignition-dist/nginx-ignition.service /etc/systemd/system/
```

### Fix the permissions
```bash
sudo chown -R root:root /opt/nginx-ignition
sudo chmod 0755 /opt/nginx-ignition/nginx-ignition
sudo chmod 0644 /opt/nginx-ignition/nginx-ignition.properties
```

If you prefer to run as a non-root user, remember to also change the user in /etc/systemd/system/nginx-ignition.service, 
set appropriate ownership of /opt/nginx-ignition and /tmp/nginx-ignition, and grant the needed permissions for the user
to interact with the nginx binary and alike.

### Configure nginx ignition
Open /opt/nginx-ignition/nginx-ignition.properties and adjust values as needed. Please check the documentation 
available on GitHub for full details, default values and recommendations.

### Register and start the service (using systemd)
Reload systemd, enable, and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable nginx-ignition
sudo systemctl start nginx-ignition
```

Check status and logs:
```bash
systemctl status nginx-ignition --no-pager
journalctl -u nginx-ignition -f
```

### Open firewall (if applicable)
- UFW: `sudo ufw allow 8090/tcp`
- firewalld: `sudo firewall-cmd --add-port=8090/tcp --permanent && sudo firewall-cmd --reload`

### Configure SELinux (if enforcing)
If SELinux blocks port binding or file access, you may need to allow the port and label the directories accordingly

```bash
sudo semanage port -a -t http_port_t -p tcp 8090 || true
sudo chcon -R -t usr_t /opt/nginx-ignition
```

### Open nginx in your browser and start using
Open http://<your-host>:8090 in your browser. The nginx ignition UI will guide you through the remaining setup steps.

## Upgrading
1. Stop service: `sudo systemctl stop nginx-ignition`
2. Replace /opt/nginx-ignition/nginx-ignition and any updated assets (frontend, migrations)
3. Start service: `sudo systemctl start nginx-ignition`
