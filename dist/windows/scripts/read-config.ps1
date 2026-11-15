# Transcodes the UTF-8 configuration file into UTF-16LE so the installer can
# read it with FileReadUTF16LE without relying on the Windows ANSI code page.

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ConfigPath,

    [Parameter(Mandatory = $true)]
    [string]$OutputPath
)

$ErrorActionPreference = "Stop"

$text = [System.IO.File]::ReadAllText($ConfigPath, [System.Text.Encoding]::UTF8)
$encoding = [System.Text.UnicodeEncoding]::new($false, $false)
[System.IO.File]::WriteAllText($OutputPath, $text, $encoding)
