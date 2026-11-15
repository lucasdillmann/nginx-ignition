# Writes or updates the nginx ignition configuration file.
#
# The installer passes the configuration values through a UTF-16LE file so that
# non-ASCII characters survive the trip from the installer to PowerShell. This
# script then writes the configuration as UTF-8 without a BOM and with Unix
# line endings, which is what the application expects.
#
# On a fresh installation the template is expanded. On an upgrade only the keys
# managed by the installer are replaced, so comments and custom keys added by the
# administrator are preserved.

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ConfigPath,

    [Parameter(Mandatory = $true)]
    [string]$TemplatePath,

    [Parameter(Mandatory = $true)]
    [string]$ValuesPath
)

$ErrorActionPreference = "Stop"

function Write-Utf8Lines {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path,

        [Parameter(Mandatory = $true)]
        [AllowEmptyCollection()]
        [AllowEmptyString()]
        [string[]]$Lines
    )

    $text = ($Lines -join "`n") + "`n"
    $encoding = [System.Text.UTF8Encoding]::new($false)
    [System.IO.File]::WriteAllText($Path, $text, $encoding)
}

# Value order, written by the installer:
#   0  nginx binary path       7  database host
#   1  nginx config path       8  database port
#   2  VPN config path         9  database name
#   3  server address         10  database schema
#   4  server port            11  database username
#   5  data path              12  database password
#   6  database driver        13  database SSL mode
$values = [System.IO.File]::ReadAllLines($ValuesPath, [System.Text.Encoding]::Unicode)
if ($values.Count -lt 14) {
    throw "Expected 14 configuration values, found $($values.Count)."
}

$driver = $values[6].Trim().ToLowerInvariant()
if ($driver -ne "sqlite" -and $driver -ne "postgres") {
    throw "Unsupported database driver: $driver"
}

$managed = [ordered]@{
    "nginx-ignition.nginx.binary-path" = $values[0]
    "nginx-ignition.nginx.config-path" = $values[1]
    "nginx-ignition.vpn.config-path"   = $values[2]
    "nginx-ignition.server.address"    = $values[3]
    "nginx-ignition.server.port"       = $values[4]
    "nginx-ignition.database.driver"   = $driver
}

if ($driver -eq "sqlite") {
    $managed["nginx-ignition.database.data-path"] = $values[5]
} else {
    $managed["nginx-ignition.database.host"] = $values[7]
    $managed["nginx-ignition.database.port"] = $values[8]
    $managed["nginx-ignition.database.name"] = $values[9]
    $managed["nginx-ignition.database.schema"] = $values[10]
    $managed["nginx-ignition.database.username"] = $values[11]
    $managed["nginx-ignition.database.password"] = $values[12]
    $managed["nginx-ignition.database.ssl-mode"] = $values[13]
}

if (-not (Test-Path -LiteralPath $ConfigPath)) {
    $template = [System.IO.File]::ReadAllText($TemplatePath, [System.Text.Encoding]::UTF8)
    $replacements = @{
        "@NGINX_BINARY_PATH@" = $values[0]
        "@NGINX_CONFIG_PATH@" = $values[1]
        "@VPN_CONFIG_PATH@"   = $values[2]
        "@SERVER_ADDRESS@"    = $values[3]
        "@SERVER_PORT@"       = $values[4]
        "@DATA_PATH@"         = $values[5]
        "@DB_HOST@"           = $values[7]
        "@DB_PORT@"           = $values[8]
        "@DB_NAME@"           = $values[9]
        "@DB_SCHEMA@"         = $values[10]
        "@DB_USERNAME@"       = $values[11]
        "@DB_PASSWORD@"       = $values[12]
        "@DB_SSL_MODE@"       = $values[13]
    }

    foreach ($placeholder in $replacements.Keys) {
        $template = $template.Replace($placeholder, $replacements[$placeholder])
    }

    $lines = $template.TrimEnd([char[]]"`r`n") -split "`r?`n"
    Write-Utf8Lines -Path $ConfigPath -Lines $lines
    exit 0
}

# Existing configuration: only touch the keys managed by the installer.
$lines = [System.IO.File]::ReadAllLines($ConfigPath, [System.Text.Encoding]::UTF8)
$output = New-Object System.Collections.Generic.List[string]
$seen = @{}

foreach ($line in $lines) {
    if ($line -match "^\s*([^#;=\s][^=]*?)\s*=") {
        $key = $matches[1]
        if ($managed.Contains($key)) {
            $output.Add("$key=$($managed[$key])")
            $seen[$key] = $true
            continue
        }
    }

    $output.Add($line)
}

foreach ($key in $managed.Keys) {
    if (-not $seen.ContainsKey($key)) {
        $output.Add("$key=$($managed[$key])")
    }
}

Write-Utf8Lines -Path $ConfigPath -Lines $output
