; Reading, generating and updating the nginx ignition properties file.

; Reads the keys managed by the installer from the existing configuration.
Function LoadExistingConfig
  Push $0
  Push $1
  Push $2
  Push $3
  Push $4
  Push $9

  StrCpy $PreviousDriver ""

  ${If} $ConfigLoadPath == ""
    Goto done
  ${EndIf}
  ${IfNot} ${FileExists} "$ConfigLoadPath"
    Goto done
  ${EndIf}

  ; The file is UTF-8. Transcode it to UTF-16LE first so FileReadUTF16LE does
  ; not decode non-ASCII paths or passwords through the ANSI code page.
  Delete "$PLUGINSDIR\config-read.txt"
  nsExec::ExecToStack `powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$PLUGINSDIR\read-config.ps1" -ConfigPath "$ConfigLoadPath" -OutputPath "$PLUGINSDIR\config-read.txt"`
  Pop $0
  Pop $1
  ${If} $0 != "0"
    Goto done
  ${EndIf}

  ClearErrors
  FileOpen $9 "$PLUGINSDIR\config-read.txt" r
  IfErrors done

  read:
    FileReadUTF16LE $9 $0
    IfErrors close

    ${StrTrimNewLines} $0 "$0"

    StrCpy $1 "$0" 1
    StrCmp $1 "#" read
    StrCmp $1 ";" read
    StrCmp $0 "" read

    ${StrLoc} $1 "$0" "=" ">"
    StrCmp $1 "" read

    StrCpy $2 "$0" $1
    IntOp $3 $1 + 1
    StrCpy $4 "$0" "" $3

    ${If} $2 == "nginx-ignition.nginx.binary-path"
      StrCpy $NginxBinaryPath "$4"
    ${ElseIf} $2 == "nginx-ignition.nginx.config-path"
      StrCpy $NginxConfigPath "$4"
    ${ElseIf} $2 == "nginx-ignition.vpn.config-path"
      StrCpy $VpnConfigPath "$4"
    ${ElseIf} $2 == "nginx-ignition.server.address"
      StrCpy $ServerAddress "$4"
    ${ElseIf} $2 == "nginx-ignition.server.port"
      StrCpy $ServerPort "$4"
    ${ElseIf} $2 == "nginx-ignition.database.driver"
      ${StrCase} $4 "$4" "L"
      ${If} $4 == "postgres"
        StrCpy $DatabaseDriver "postgres"
      ${Else}
        StrCpy $DatabaseDriver "sqlite"
      ${EndIf}
      StrCpy $PreviousDriver "$DatabaseDriver"
    ${ElseIf} $2 == "nginx-ignition.database.data-path"
      StrCpy $DataDir "$4"
    ${ElseIf} $2 == "nginx-ignition.database.host"
      StrCpy $DbHost "$4"
    ${ElseIf} $2 == "nginx-ignition.database.port"
      StrCpy $DbPort "$4"
    ${ElseIf} $2 == "nginx-ignition.database.name"
      StrCpy $DbName "$4"
    ${ElseIf} $2 == "nginx-ignition.database.schema"
      StrCpy $DbSchema "$4"
    ${ElseIf} $2 == "nginx-ignition.database.username"
      StrCpy $DbUsername "$4"
    ${ElseIf} $2 == "nginx-ignition.database.password"
      StrCpy $DbPassword "$4"
      StrCpy $ExistingDbPassword "$4"
    ${ElseIf} $2 == "nginx-ignition.database.ssl-mode"
      StrCpy $DbSslMode "$4"
    ${EndIf}

    Goto read

  close:
    FileClose $9
    Delete "$PLUGINSDIR\config-read.txt"

  done:
  Pop $9
  Pop $4
  Pop $3
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

Function EnsureDirectories
  CreateDirectory "$APPDATA_DIR"
  CreateDirectory "$INSTDIR"

  ${If} $DataDir != ""
    CreateDirectory "$DataDir"
  ${EndIf}
  ${If} $NginxConfigPath != ""
    CreateDirectory "$NginxConfigPath"
  ${EndIf}
  ${If} $VpnConfigPath != ""
    CreateDirectory "$VpnConfigPath"
  ${EndIf}
FunctionEnd

; Writes the managed configuration values. On an upgrade without explicit
; configuration arguments the existing file is left untouched.
Function WriteConfiguration
  ${If} $ConfigImport != ""
    ClearErrors
    CopyFiles /SILENT "$ConfigImport" "$CONFIG_FILE"
    ${If} ${Errors}
      MessageBox MB_ICONSTOP|MB_OK "$(NI_CONFIG_IMPORT_FAILED)"
      Abort
    ${EndIf}
    Return
  ${EndIf}

  ${If} $IsUpgrade == "1"
  ${AndIf} $ConfigTouched != "1"
    ${If} ${FileExists} "$CONFIG_FILE"
      Return
    ${EndIf}
  ${ElseIf} $HasExistingConfig == "1"
  ${AndIf} $ConfigTouched != "1"
  ${AndIf} ${Silent}
    ${If} ${FileExists} "$CONFIG_FILE"
      Return
    ${EndIf}
  ${EndIf}

  DetailPrint "$(NI_DETAIL_CONFIG)"

  ClearErrors
  FileOpen $9 "$PLUGINSDIR\config-values.txt" w
  ${If} ${Errors}
    MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_CONFIG)"
    Abort
  ${EndIf}

  ; The values file is UTF-16LE so that non-ASCII paths and passwords survive.
  FileWriteUTF16LE $9 "$NginxBinaryPath$\r$\n"
  FileWriteUTF16LE $9 "$NginxConfigPath$\r$\n"
  FileWriteUTF16LE $9 "$VpnConfigPath$\r$\n"
  FileWriteUTF16LE $9 "$ServerAddress$\r$\n"
  FileWriteUTF16LE $9 "$ServerPort$\r$\n"
  FileWriteUTF16LE $9 "$DataDir$\r$\n"
  FileWriteUTF16LE $9 "$DatabaseDriver$\r$\n"
  FileWriteUTF16LE $9 "$DbHost$\r$\n"
  FileWriteUTF16LE $9 "$DbPort$\r$\n"
  FileWriteUTF16LE $9 "$DbName$\r$\n"
  FileWriteUTF16LE $9 "$DbSchema$\r$\n"
  FileWriteUTF16LE $9 "$DbUsername$\r$\n"
  FileWriteUTF16LE $9 "$DbPassword$\r$\n"
  FileWriteUTF16LE $9 "$DbSslMode$\r$\n"
  FileClose $9

  ${If} $DatabaseDriver == "postgres"
    StrCpy $0 "$PLUGINSDIR\nginx-ignition.postgres.properties"
  ${Else}
    StrCpy $0 "$PLUGINSDIR\nginx-ignition.sqlite.properties"
  ${EndIf}

  nsExec::ExecToStack `powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$PLUGINSDIR\update-config.ps1" -ConfigPath "$CONFIG_FILE" -TemplatePath "$0" -ValuesPath "$PLUGINSDIR\config-values.txt"`
  Pop $1
  Pop $2

  ${If} $1 != "0"
    MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_CONFIG)$\r$\n$\r$\n$2"
    Abort
  ${EndIf}

  Delete "$PLUGINSDIR\config-values.txt"
FunctionEnd
