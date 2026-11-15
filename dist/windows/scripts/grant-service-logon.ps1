# Grants the "Log on as a service" right (SeServiceLogonRight) to an account.
#
# The nginx ignition installer uses this for custom service accounts. Built-in
# accounts (LocalSystem, LocalService and NetworkService) already have the
# right and do not need this script.
#
# Usage:
#   powershell.exe -NoProfile -ExecutionPolicy Bypass -File grant-service-logon.ps1 -Account "DOMAIN\User"
#
# The script exits with code 0 on success and 1 on failure.

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$Account
)

$ErrorActionPreference = "Stop"

$source = @"
using System;
using System.Runtime.InteropServices;

public static class LsaHelper
{
    [StructLayout(LayoutKind.Sequential)]
    public struct LSA_UNICODE_STRING
    {
        public ushort Length;
        public ushort MaximumLength;
        public IntPtr Buffer;
    }

    [StructLayout(LayoutKind.Sequential)]
    public struct LSA_OBJECT_ATTRIBUTES
    {
        public int Length;
        public IntPtr RootDirectory;
        public IntPtr ObjectName;
        public uint Attributes;
        public IntPtr SecurityDescriptor;
        public IntPtr SecurityQualityOfService;
    }

    [DllImport("advapi32.dll", SetLastError = true)]
    public static extern uint LsaOpenPolicy(
        IntPtr SystemName,
        ref LSA_OBJECT_ATTRIBUTES ObjectAttributes,
        uint AccessMask,
        out IntPtr PolicyHandle);

    [DllImport("advapi32.dll", SetLastError = true)]
    public static extern uint LsaAddAccountRights(
        IntPtr PolicyHandle,
        byte[] AccountSid,
        LSA_UNICODE_STRING[] UserRights,
        uint CountOfRights);

    [DllImport("advapi32.dll", SetLastError = true)]
    public static extern uint LsaClose(IntPtr ObjectHandle);
}
"@

try {
    Add-Type -TypeDefinition $source -ErrorAction Stop

    $sid = (New-Object -TypeName System.Security.Principal.NTAccount -ArgumentList $Account).Translate(
        [System.Security.Principal.SecurityIdentifier])
    $sidBytes = New-Object byte[] $sid.BinaryLength
    $sid.GetBinaryForm($sidBytes, 0)

    $right = "SeServiceLogonRight"
    $rightBytes = [System.Text.Encoding]::Unicode.GetBytes($right)
    $rightBuffer = [System.Runtime.InteropServices.Marshal]::AllocHGlobal($rightBytes.Length + 2)

    try {
        [System.Runtime.InteropServices.Marshal]::Copy($rightBytes, 0, $rightBuffer, $rightBytes.Length)
        [System.Runtime.InteropServices.Marshal]::WriteByte($rightBuffer, $rightBytes.Length, 0)
        [System.Runtime.InteropServices.Marshal]::WriteByte($rightBuffer, $rightBytes.Length + 1, 0)

        $lsaRight = New-Object -TypeName "LsaHelper+LSA_UNICODE_STRING"
        $lsaRight.Length = [ushort]$rightBytes.Length
        $lsaRight.MaximumLength = [ushort]($rightBytes.Length + 2)
        $lsaRight.Buffer = $rightBuffer

        $attributes = New-Object -TypeName "LsaHelper+LSA_OBJECT_ATTRIBUTES"
        $attributes.Length = [System.Runtime.InteropServices.Marshal]::SizeOf($attributes)

        $policyHandle = [IntPtr]::Zero
        # POLICY_ALL_ACCESS is required by LsaAddAccountRights.
        $status = [LsaHelper]::LsaOpenPolicy(
            [IntPtr]::Zero, [ref]$attributes, 0x000F0FFF, [ref]$policyHandle)
        if ($status -ne 0) {
            throw ("LsaOpenPolicy failed with status 0x{0:X8}" -f $status)
        }

        try {
            $status = [LsaHelper]::LsaAddAccountRights($policyHandle, $sidBytes, @($lsaRight), 1)
            if ($status -ne 0) {
                throw ("LsaAddAccountRights failed with status 0x{0:X8}" -f $status)
            }
        }
        finally {
            [void][LsaHelper]::LsaClose($policyHandle)
        }
    }
    finally {
        [System.Runtime.InteropServices.Marshal]::FreeHGlobal($rightBuffer)
    }

    Write-Output "Granted SeServiceLogonRight to $Account"
    exit 0
}
catch {
    Write-Error $_.Exception.Message
    exit 1
}
