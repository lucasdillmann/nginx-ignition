# Windows installation and upgrade instructions

These instructions explain how to install and run nginx-ignition on a Windows host. The distribution .zip you downloaded 
already includes:

- nginx-ignition binary (nginx-ignition.exe)
- frontend assets (UI)
- database migrations
- default configuration file, nginx-ignition.properties

## Requirements
- Windows (amd64 or arm64)
- nginx installed and accessible (can be changed in the properties file)
- nginx modules mod-http-js, mod-http-lua and mod-stream installed
- Administrator privileges for installation and running (if binding to privileged ports)
- Port 8090 available (default UI/API port; also configurable)

## New installations
### Unpack the distribution
Extract the ZIP archive you downloaded (windows-amd64 or windows-arm64) to a directory of your choice, 
`C:\nginx-ignition` for example.

### Configure nginx ignition
Open `nginx-ignition.properties` and adjust values as needed. Make sure to set the `nginx-ignition.nginx.binary-path`
to the correct location of your `nginx.exe`.

nginx ignition can also be configured using environment variables. The keys are the same as in the properties file, 
but with uppercase names and dots/dashes replaced by underscores. For example, the `nginx-ignition.nginx.binary-path` 
property can be set using the `NGINX_IGNITION_NGINX_BINARY_PATH` environment variable.

### Running the application
To start nginx ignition, open a Command Prompt or PowerShell as Administrator and run:

```powershell
cd C:\nginx-ignition
.\nginx-ignition.exe
```

### Open firewall (if applicable)
If you have Windows Firewall enabled, you'll need to allow traffic on port 8090:

```powershell
New-NetFirewallRule -DisplayName "nginx-ignition" -Direction Inbound -LocalPort 8090 -Protocol TCP -Action Allow
```

### Open nginx in your browser and start using
Open http://localhost:8090 in your browser. The nginx ignition UI will guide you through the remaining setup steps.

## Upgrading
1. Stop the running `nginx-ignition.exe` process.
2. Replace `nginx-ignition.exe` and any updated assets (`frontend/`, `migrations/`).
3. Start the application again.
