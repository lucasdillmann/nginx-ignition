# Waits for a Windows service to reach the requested state.
#
# The installer uses this instead of parsing the localized output of sc.exe.
# Exit codes:
#   0  the service reached the requested state
#   1  the service did not reach the requested state before the timeout
#   2  the service does not exist

[CmdletBinding()]
param(
    [Parameter(Mandatory = $false)]
    [string]$Name = "nginx-ignition",

    [Parameter(Mandatory = $false)]
    [ValidateSet("Running", "Stopped")]
    [string]$Status = "Running",

    [Parameter(Mandatory = $false)]
    [int]$TimeoutSeconds = 30,

    # When set, check the current state without waiting for a transition.
    [Parameter(Mandatory = $false)]
    [switch]$Check
)

$ErrorActionPreference = "Stop"

$service = Get-Service -Name $Name -ErrorAction SilentlyContinue
if ($null -eq $service) {
    exit 2
}

if ($Check) {
    if ($service.Status -eq $Status) {
        exit 0
    }

    exit 1
}

try {
    $service.WaitForStatus([System.ServiceProcess.ServiceControllerStatus]$Status,
        [TimeSpan]::FromSeconds($TimeoutSeconds))
}
catch {
    exit 1
}

if ($service.Status -eq $Status) {
    exit 0
}

exit 1
