# Windows installer

The Windows distribution is available as an NSIS installer (`nginx-ignition-<version>.windows-<arch>-setup.exe`) and
as a portable ZIP archive. This document describes the installer. For user-facing installation steps, see
[`dist/windows/instructions.md`](../dist/windows/instructions.md).

## What the installer does

- Copies the single-binary `nginx-ignition.exe` to `C:\Program Files\nginx-ignition`.
- Writes a commented properties file to `C:\ProgramData\nginx-ignition\nginx-ignition.properties`.
- Creates the data, managed nginx configuration and VPN directories under `C:\ProgramData\nginx-ignition`.
- Registers the `nginx-ignition` Windows service with the selected account and startup type.
- Creates Windows Firewall rules for the web interface (domain/private profiles) and, optionally, ports 80 and 443.
- Registers `nginx-ignition` as an Event Log source so service logs are visible in Event Viewer.
- Grants `SYSTEM`, `Administrators` and the service account access to the configuration and data directories.
- Adds Start Menu and, optionally, desktop shortcuts that open the web interface in the default browser, plus an
  entry in **Apps & features**.
- Preserves the configuration and data during upgrades and offers to keep them during uninstallation.

The installer requires Windows 10 or Windows Server 2016 or later. Separate installers are produced for amd64 and
arm64; each refuses to install on an incompatible architecture.

## Installation modes

### Express

Express is intended for the common case and requires no decisions beyond the installation directory:

- Nginx is auto-detected (`%ProgramFiles%\nginx\nginx.exe`, `C:\nginx\nginx.exe`, then `PATH`)
- Listen address `127.0.0.1` (localhost only), port `8090`
- Embedded SQLite database in `C:\ProgramData\nginx-ignition\data` (recommended for tests; use PostgreSQL
  for production)
- Service runs as `LocalSystem`, starts automatically
- Firewall rule for the web interface
- Start Menu shortcut

### Advanced

Advanced exposes every installer-managed setting on separate wizard pages:

- **nginx**: `nginx.exe` path, managed nginx configuration directory, VPN configuration directory
- **Server and data**: listen address, port, data directory
- **Database**: embedded SQLite (recommended for tests) or PostgreSQL (recommended for production), including
  host, port, name, schema, user, password and SSL mode
- **Service**: `LocalSystem`, `Network Service`, `Local Service` or a custom account, plus automatic/manual startup.
  Choosing a custom account opens a follow-up page for the account name and password.
- **Firewall and shortcuts**: web interface, HTTP 80 and HTTPS 443 rules, desktop and Start Menu shortcuts

## Managed configuration keys

The installer writes the following keys and preserves every other line of an existing properties file:

| Key                                          | Notes                                        |
|:---------------------------------------------|:---------------------------------------------|
| `nginx-ignition.nginx.binary-path`           | Absolute path to `nginx.exe`                 |
| `nginx-ignition.nginx.config-path`           | Directory for nginx ignition-managed configs |
| `nginx-ignition.vpn.config-path`             | Directory for VPN configuration              |
| `nginx-ignition.server.address`              | Listen address                               |
| `nginx-ignition.server.port`                 | UI/API port                                  |
| `nginx-ignition.database.driver`             | `sqlite` or `postgres`                       |
| `nginx-ignition.database.data-path`          | SQLite data directory                        |
| `nginx-ignition.database.host`               | PostgreSQL only                              |
| `nginx-ignition.database.port`               | PostgreSQL only                              |
| `nginx-ignition.database.name`               | PostgreSQL only                              |
| `nginx-ignition.database.schema`             | PostgreSQL only                              |
| `nginx-ignition.database.username`           | PostgreSQL only                              |
| `nginx-ignition.database.password`           | PostgreSQL only                              |
| `nginx-ignition.database.ssl-mode`           | `disable` or `require`                       |

The generated file is UTF-8 without a BOM and uses Unix line endings, which is what the application expects.

## Silent installation

Run the installer with `/S`. The following arguments are supported:

| Argument                    | Meaning                                                       |
|:----------------------------|:--------------------------------------------------------------|
| `/NGINX=<path>`             | Path to `nginx.exe`                                           |
| `/NGINXCONFIG=<path>`       | Managed nginx configuration directory                         |
| `/VPNCONFIG=<path>`         | VPN configuration directory                                   |
| `/ADDRESS=<address>`        | Listen address (default `127.0.0.1`)                          |
| `/PORT=<port>`              | UI/API port (default `8090`)                                  |
| `/DATADIR=<path>`           | Data directory (default `C:\ProgramData\nginx-ignition\data`) |
| `/DBDRIVER=<driver>`        | `sqlite` or `postgres`                                        |
| `/DBHOST=<host>`            | PostgreSQL host                                               |
| `/DBPORT=<port>`            | PostgreSQL port (default `5432`)                              |
| `/DBNAME=<name>`            | PostgreSQL database name                                      |
| `/DBSCHEMA=<schema>`        | PostgreSQL schema (default `public`)                          |
| `/DBUSER=<user>`            | PostgreSQL username                                           |
| `/DBPASS=<password>`        | PostgreSQL password (required)                                |
| `/DBSSLMODE=<mode>`         | `disable` or `require` (default `require`)                    |
| `/SERVICEACCOUNT=<kind>`    | `localsystem`, `networkservice`, `localservice` or `custom`   |
| `/SERVICEUSER=<account>`    | Account name for a custom service account                     |
| `/SERVICEPASS=<password>`   | Password for a custom service account                         |
| `/SERVICESTART=<type>`      | `auto` or `manual`                                            |
| `/START=<bool>`             | Start the service after installation                          |
| `/FIREWALL=<bool>`          | Allow the web interface port on domain/private profiles       |
| `/FIREWALLHTTP=<bool>`      | Allow TCP port 80                                             |
| `/FIREWALLHTTPS=<bool>`     | Allow TCP port 443                                            |
| `/DESKTOPSHORTCUT=<bool>`   | Create a desktop shortcut                                     |
| `/STARTMENUSHORTCUT=<bool>` | Create a Start Menu shortcut                                  |
| `/OPENUI=<bool>`            | Open the web interface after a successful silent installation |
| `/CONFIG=<path>`            | Copy this properties file instead of generating one           |
| `/INSTALLTYPE=<type>`       | `express` or `advanced` (only affects defaults)               |
| `/LANG=<language id>`       | Select the installer language (NSIS built-in)                 |
| `/D=<path>`                 | Installation directory; must be the last argument             |

Boolean arguments accept `1`/`0`, `true`/`false`, `yes`/`no` and `on`/`off`.

### Examples

Express, silent:

```powershell
nginx-ignition-setup.exe /S
```

PostgreSQL with a custom port:

```powershell
nginx-ignition-setup.exe /S `
  /DBDRIVER=postgres `
  /DBHOST=db.example.com `
  /DBPORT=5432 `
  /DBNAME=nginx_ignition `
  /DBSCHEMA=public `
  /DBUSER=nginx_ignition `
  /DBPASS="secret" `
  /DBSSLMODE=require `
  /PORT=8090
```

Reusing a prepared configuration:

```powershell
nginx-ignition-setup.exe /S /CONFIG="\\fileserver\deploy\nginx-ignition.properties"
```

The imported file is used as-is, so it must contain absolute paths (the service runs with `C:\Windows\System32` as its
working directory). The installer still creates and grants access to the default directories; if the imported file
points elsewhere, make sure the service account can read the file and write to the configured data, nginx and VPN
directories.

## Upgrade behavior

An installation is considered an upgrade when the uninstall registry key exists. During an upgrade the installer:

- Stops the service and remembers whether it was running,
- Skips the installation-type and configuration pages,
- Keeps the service object, account, password and startup type,
- Updates only the service binary path and description,
- Leaves the properties file byte-for-byte untouched unless configuration arguments or `/CONFIG` are provided,
- Copies the new binary over the old one,
- Restarts the service if it was running before.

Because the service object is preserved, upgrades do not require the password of a custom service account.

## Uninstallation

The uninstaller stops and removes the service, firewall rules, Event Log source, shortcuts and registry entries, then
removes the application files.

Configuration and data in `C:\ProgramData\nginx-ignition` are kept by default. The interactive uninstaller asks before
deleting them. For silent uninstalls, use `/PURGE=1` to delete them:

```powershell
uninstall.exe /S /PURGE=1
```

Custom data directories outside `C:\ProgramData\nginx-ignition` are never deleted.

## Troubleshooting

### Windows SmartScreen warns about the installer

The installer is not Authenticode-signed yet. If SmartScreen shows a warning, choose **More info** and then
**Run anyway**, or verify the SHA-256 checksum of the download against the release page.

### The installer refuses the architecture

Use the arm64 installer on Windows on ARM64 and the amd64 installer everywhere else.

### The service does not start

Check Event Viewer → Windows Logs → Application, filtered by the `nginx-ignition` source.

### A custom service account fails to log on

Grant the account the "Log on as a service" right in `secpol.msc` and ensure the password is correct. The installer
attempts to grant the right automatically but cannot do so in every environment.

### The installer cannot write the configuration file

Check the permissions of `C:\ProgramData\nginx-ignition`. Antivirus or group policy may block writes there.
