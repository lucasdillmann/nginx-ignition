# Windows installation and upgrade instructions

nginx ignition is distributed for Windows in two forms:

- **Installer (`nginx-ignition-<version>.windows-<arch>-setup.exe`)**: recommended. Installs the application under
  `C:\Program Files\nginx-ignition`, keeps configuration and data under `C:\ProgramData\nginx-ignition`, registers the
  Windows service, firewall rules and the Event Log source, and adds Start Menu/desktop shortcuts.
- **ZIP archive (`nginx-ignition-<version>.windows-<arch>.zip`)**: the portable layout. You copy the files wherever you
  want and manage the service, firewall and permissions yourself.

Both packages contain the same single-binary nginx ignition build. nginx is **not** included and must be installed
separately.

## Requirements

- Windows 10 or Windows Server 2016 or later (amd64 or arm64)
- nginx installed and accessible (can be changed in the properties file or the installer)
- Administrator privileges for installation and, unless you pick an unprivileged port, for running
- Port 8090 available by default (the UI/API port is configurable)

## Installing with the installer

1. Run the `setup.exe` for your architecture. The installer refuses to run on the wrong architecture and requires
   Windows 10/Server 2016 or later.
2. Choose the display language. It is remembered for upgrades and uninstalls. The Windows binaries and installers are
   Authenticode-signed. Free code signing provided by [SignPath.io](https://about.signpath.io), certificate by
   [SignPath Foundation](https://signpath.org). See the
   [code signing policy](https://github.com/lucasdillmann/nginx-ignition/blob/main/CODE_SIGNING.md).
3. Choose **Express** or **Advanced**:
   - **Express** uses the detected nginx installation, listens on `127.0.0.1:8090`, stores data in
     `C:\ProgramData\nginx-ignition\data`, uses the embedded SQLite database (recommended for tests; use PostgreSQL for
     production), runs the service as `LocalSystem` and opens the web interface when setup finishes.
   - **Advanced** additionally lets you configure the nginx paths, listen address/port, data directory, SQLite or
     PostgreSQL, the service account and startup type, firewall rules and shortcuts.
4. Complete the wizard. The installer starts the service and adds a Start Menu entry.

To open nginx ignition later, use the Start Menu (or desktop) shortcut. It opens the web interface in your default
browser, using the address and port configured during installation.

### File locations

| Content                          | Location                                                   |
|:---------------------------------|:-----------------------------------------------------------|
| Application and uninstaller      | `C:\Program Files\nginx-ignition`                          |
| Configuration                    | `C:\ProgramData\nginx-ignition\nginx-ignition.properties`  |
| SQLite database and data         | `C:\ProgramData\nginx-ignition\data` (default)             |
| Managed nginx configuration      | `C:\ProgramData\nginx-ignition\nginx` (default)            |
| VPN configuration                | `C:\ProgramData\nginx-ignition\vpn` (default)              |

The installer removes the inherited access from `C:\ProgramData\nginx-ignition` (the properties file contains
database credentials), grants `SYSTEM` and `Administrators` full control, and grants the service account the
permissions it needs. Program Files keeps its standard read/execute permissions.

### Managing the service

```powershell
Get-Service nginx-ignition
Stop-Service nginx-ignition
Start-Service nginx-ignition
Restart-Service nginx-ignition
```

When running as a service, logs are written to the Windows Event Log under the `nginx-ignition` source.

### Silent installation

All wizard options are available as command-line arguments, which makes the installer suitable for automation. Use
`/S` for a silent install. `makensis`/NSIS also accepts `/LANG=<language id>` to select the language, and `/D=<path>`
must be the **last** argument if you want to change the installation directory.

```powershell
nginx-ignition-setup.exe /S `
  /NGINX="C:\nginx\nginx.exe" `
  /NGINXCONFIG="C:\ProgramData\nginx-ignition\nginx" `
  /VPNCONFIG="C:\ProgramData\nginx-ignition\vpn" `
  /ADDRESS=127.0.0.1 `
  /PORT=8090 `
  /DATADIR="C:\ProgramData\nginx-ignition\data" `
  /DBDRIVER=sqlite `
  /SERVICEACCOUNT=localsystem `
  /SERVICESTART=auto `
  /START=1 `
  /FIREWALL=1 `
  /FIREWALLHTTP=0 `
  /FIREWALLHTTPS=0 `
  /DESKTOPSHORTCUT=0 `
  /STARTMENUSHORTCUT=1
```

For PostgreSQL, set `/DBDRIVER=postgres` and provide `/DBHOST`, `/DBPORT`, `/DBNAME`, `/DBSCHEMA`, `/DBUSER`,
`/DBPASS` and `/DBSSLMODE` (`disable` or `require`). A non-empty password is required.

To reuse a complete, pre-prepared properties file instead of generating one, pass `/CONFIG=<path>`:

```powershell
setup.exe /S /CONFIG="\\fileserver\deploy\nginx-ignition.properties"
```

`/CONFIG` takes precedence over the individual configuration arguments and leaves the provided file untouched. The
installer uses its listen address, port and configured data/nginx/VPN directories for the firewall rules, shortcuts
and service-account permissions. Use absolute paths in the imported file: the service runs with
`C:\Windows\System32` as its working directory.

Boolean arguments accept `1`/`0`, `true`/`false`, `yes`/`no` and `on`/`off`. `/START` controls whether the service is
started after installation. `/OPENUI` is `0` by default in silent mode; set `/OPENUI=1` to open the web interface in
the default browser after a successful silent installation.

### Custom service accounts

The installer can run the service as `LocalSystem` (default), `Network Service`, `Local Service` or a custom account:

```powershell
setup.exe /S /SERVICEACCOUNT=custom /SERVICEUSER="DOMAIN\svc-nginx-ignition" /SERVICEPASS="secret"
```

For custom accounts the installer grants the "Log on as a service" right on a best-effort basis. If that fails, grant
it manually with the Local Security Policy snap-in (`secpol.msc` → Local Policies → User Rights Assignment → Log on
as a service) and check the Windows Event Log if the service does not start. Passwords containing double quotes are not
supported.

### Upgrading

Run the setup for the new version. The installer:

- Stops the running service and preserves its credentials and startup type,
- Replaces the application files in the existing installation directory,
- Keeps `nginx-ignition.properties`, the SQLite database and any other data unchanged,
- Restarts the service if it was running before the upgrade.

Passing configuration arguments (or `/CONFIG`) during an upgrade updates only the keys managed by the installer; comments
and custom keys you added to the properties file are preserved. To change settings interactively, edit the properties
file or run the installer with the desired arguments.

### Uninstalling

Use **Settings → Apps → Installed apps → nginx ignition → Uninstall**, or run `uninstall.exe` from the installation
directory. The uninstaller stops and removes the service, firewall rules, shortcuts and Event Log source.

By default, configuration and data in `C:\ProgramData\nginx-ignition` are **kept** so a future installation can reuse
them. The interactive uninstaller asks before deleting them. For a silent uninstall, pass `/PURGE=1` to delete them:

```powershell
uninstall.exe /S /PURGE=1
```

A silent uninstall without `/PURGE=1` keeps the data. Custom data directories configured outside
`C:\ProgramData\nginx-ignition` are never deleted by the uninstaller.

If you reinstall after uninstalling with data kept, setup detects the existing `nginx-ignition.properties` and reuses
it. A silent reinstall without configuration arguments leaves it untouched; an interactive reinstall pre-fills the
configuration pages from it.

## Installing from the ZIP archive

### Unpack the distribution

Extract the ZIP archive to a directory of your choice, `C:\nginx-ignition` for example.

### Configure nginx ignition

Open `nginx-ignition.properties` and adjust values as needed. Make sure to set
`nginx-ignition.nginx.binary-path` to the correct location of your `nginx.exe`.

nginx ignition can also be configured using environment variables. The keys are the same as in the properties file,
but with uppercase names and dots/dashes replaced by underscores. For example, the
`nginx-ignition.nginx.binary-path` property can be set using the `NGINX_IGNITION_NGINX_BINARY_PATH` environment
variable.

### Running the application

To start nginx ignition manually, open a Command Prompt or PowerShell as Administrator and run:

```powershell
cd C:\nginx-ignition
.\nginx-ignition.exe
```

### Installing as a background service

nginx ignition supports the Windows Service Control Manager natively. Open a Command Prompt or PowerShell as
Administrator and register the service with the exact name `nginx-ignition`:

```powershell
sc.exe create nginx-ignition binPath= '"C:\nginx-ignition\nginx-ignition.exe" --config "C:\nginx-ignition\nginx-ignition.properties"' start= auto
sc.exe start nginx-ignition
```

When started by the Service Control Manager, nginx ignition writes its logs to the Windows Event Log under the
`nginx-ignition` source instead of stdout. The event source is registered automatically on the first service start.

### Open firewall (if applicable)

If you have Windows Firewall enabled, you'll need to allow traffic on port 8090:

```powershell
New-NetFirewallRule -DisplayName "nginx-ignition" -Direction Inbound -LocalPort 8090 -Protocol TCP -Action Allow
```

### Open nginx in your browser and start using

Open http://localhost:8090 in your browser. The nginx ignition UI will guide you through the remaining setup steps.

## Migrating from the ZIP layout to the installer

The installer does not import an existing ZIP installation automatically. To migrate:

1. Note the values in your existing `nginx-ignition.properties`, especially the nginx binary path and any PostgreSQL
   settings. Keep a copy of the file.
2. Stop the old application and remove any service you registered for it:
   ```powershell
   sc.exe stop nginx-ignition
   sc.exe delete nginx-ignition
   ```
3. Run the installer and choose **Advanced**. Configure the same nginx paths, listen address/port and database settings.
   Alternatively, pass `/CONFIG=<path-to-existing-properties>` to reuse the file as-is.
4. On the service page, clear **Start the service after installation** (or use `/START=0`), and finish the wizard.
5. Copy any data you want to keep from the old data directory into `C:\ProgramData\nginx-ignition\data`. If you used
   PostgreSQL, no data migration is needed: point the installer (or `/CONFIG`) at the same database.
6. Start the service:
   ```powershell
   sc.exe start nginx-ignition
   ```
7. Remove the old installation directory once you have verified that the new installation works.

When `/CONFIG` points at a properties file that uses directories outside `C:\ProgramData\nginx-ignition`, the
installer creates and grants access to the default directories only. Make sure the service account can read the
properties file and write to the configured data, nginx and VPN directories.
