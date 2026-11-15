; Prerequisite checks, nginx detection and command-line parsing.

Function CheckPrerequisites
  !if "${ARCH}" == "amd64"
    ${IfNot} ${RunningX64}
      MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_ARCH_X64)"
      Abort
    ${EndIf}
    ${If} ${IsNativeARM64}
      MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_ARCH_MISMATCH)"
      Abort
    ${EndIf}
  !else if "${ARCH}" == "arm64"
    ${IfNot} ${IsNativeARM64}
      MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_ARCH_ARM64)"
      Abort
    ${EndIf}
  !else
    !error "Unsupported ARCH: ${ARCH}"
  !endif

  ${IfNot} ${AtLeastWin10}
    MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_WINVER)"
    Abort
  ${EndIf}
FunctionEnd

; Input: $R0. Output: $R1 = 1 when the value is a valid TCP port.
Function ValidatePort
  Push $0
  Push $2
  Push $3
  Push $4

  StrCpy $R1 "0"
  StrLen $0 "$R0"

  ${If} $0 == 0
    Goto done
  ${EndIf}
  ${If} $0 > 5
    Goto done
  ${EndIf}

  StrCpy $2 0
  digits:
    ${If} $2 >= $0
      Goto range
    ${EndIf}
    StrCpy $3 "$R0" 1 $2
    ${StrStr} $4 "0123456789" "$3"
    ${If} $4 == ""
      Goto done
    ${EndIf}
    IntOp $2 $2 + 1
    Goto digits

  range:
    IntOp $3 "$R0" + 0
    ${If} $3 >= 1
    ${AndIf} $3 <= 65535
      StrCpy $R1 "1"
    ${EndIf}

  done:
  Pop $4
  Pop $3
  Pop $2
  Pop $0
FunctionEnd

Function DetectNginx
  ${If} $NginxBinaryPath != ""
    ${If} ${FileExists} "$NginxBinaryPath"
      Goto paths
    ${EndIf}
  ${EndIf}

  StrCpy $0 ""
  ${If} ${FileExists} "$PROGRAMFILES64\nginx\nginx.exe"
    StrCpy $0 "$PROGRAMFILES64\nginx\nginx.exe"
  ${ElseIf} ${FileExists} "$PROGRAMFILES32\nginx\nginx.exe"
    StrCpy $0 "$PROGRAMFILES32\nginx\nginx.exe"
  ${ElseIf} ${FileExists} "C:\nginx\nginx.exe"
    StrCpy $0 "C:\nginx\nginx.exe"
  ${Else}
    SearchPath $0 "nginx.exe"
  ${EndIf}

  ${If} $0 != ""
    StrCpy $NginxBinaryPath $0
  ${Else}
    StrCpy $NginxBinaryPath "$PROGRAMFILES64\nginx\nginx.exe"
  ${EndIf}

  paths:
  ${If} $NginxConfigPath == ""
    StrCpy $NginxConfigPath "$APPDATA_DIR\nginx"
  ${EndIf}
  ${If} $VpnConfigPath == ""
    StrCpy $VpnConfigPath "$APPDATA_DIR\vpn"
  ${EndIf}
FunctionEnd

Function ApplyExpressDefaults
  StrCpy $ServerAddress "127.0.0.1"
  StrCpy $ServerPort "8090"
  StrCpy $DatabaseDriver "sqlite"
  StrCpy $DataDir "$APPDATA_DIR\data"
  StrCpy $ServiceAccountKind "localsystem"
  StrCpy $ServiceAccount ""
  StrCpy $ServicePassword ""
  StrCpy $ServiceStartType "auto"
  StrCpy $ServiceStartNow "1"
  StrCpy $FirewallUI "1"
  StrCpy $FirewallHTTP "0"
  StrCpy $FirewallHTTPS "0"
  StrCpy $DesktopShortcut "0"
  StrCpy $StartMenuShortcut "1"
  Call DetectNginx
FunctionEnd

; Input: $R0. Output: $R1 = 1 for true, 0 for false.
Function ParseBoolean
  Push $0
  ${StrCase} $0 "$R0" "L"
  StrCpy $R1 "0"
  ${If} $0 == "1"
  ${OrIf} $0 == "true"
  ${OrIf} $0 == "yes"
  ${OrIf} $0 == "on"
    StrCpy $R1 "1"
  ${EndIf}
  Pop $0
FunctionEnd

Function ParseArguments
  ${GetParameters} $R9

  ${GetOptions} $R9 "/NGINX=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $NginxBinaryPath $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/NGINXCONFIG=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $NginxConfigPath $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/VPNCONFIG=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $VpnConfigPath $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/ADDRESS=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $ServerAddress $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/PORT=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ValidatePort
    ${If} $R1 != "1"
      MessageBox MB_ICONSTOP|MB_OK "$(NI_INVALID_PORT)"
      Abort
    ${EndIf}
    StrCpy $ServerPort $R2
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DATADIR=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $DataDir $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBDRIVER=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    ${StrCase} $R0 "$R0" "L"
    ${If} $R0 == "postgres"
    ${OrIf} $R0 == "postgresql"
      StrCpy $DatabaseDriver "postgres"
    ${Else}
      StrCpy $DatabaseDriver "sqlite"
    ${EndIf}
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBHOST=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $DbHost $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBPORT=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $DbPort $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBNAME=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $DbName $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBSCHEMA=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $DbSchema $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBUSER=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $DbUsername $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBPASS=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $DbPassword $R0
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/DBSSLMODE=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    ${StrCase} $R0 "$R0" "L"
    ${If} $R0 == "disable"
      StrCpy $DbSslMode "disable"
    ${Else}
      StrCpy $DbSslMode "require"
    ${EndIf}
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/SERVICEACCOUNT=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    ${StrCase} $R0 "$R0" "L"
    ${If} $R0 == "localsystem"
      StrCpy $ServiceAccountKind "localsystem"
    ${ElseIf} $R0 == "networkservice"
      StrCpy $ServiceAccountKind "networkservice"
    ${ElseIf} $R0 == "localservice"
      StrCpy $ServiceAccountKind "localservice"
    ${ElseIf} $R0 == "custom"
      StrCpy $ServiceAccountKind "custom"
    ${EndIf}
    StrCpy $ServiceAccountArgGiven "1"
  ${EndIf}

  ${GetOptions} $R9 "/SERVICEUSER=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $ServiceAccount $R0
    StrCpy $ServiceAccountKind "custom"
    StrCpy $ServiceAccountArgGiven "1"
  ${EndIf}

  ${GetOptions} $R9 "/SERVICEPASS=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $ServicePassword $R0
    StrCpy $ServiceAccountArgGiven "1"
  ${EndIf}

  ${GetOptions} $R9 "/SERVICESTART=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    ${StrCase} $R0 "$R0" "L"
    ${If} $R0 == "manual"
    ${OrIf} $R0 == "demand"
      StrCpy $ServiceStartType "manual"
    ${Else}
      StrCpy $ServiceStartType "auto"
    ${EndIf}
    StrCpy $ServiceStartArgGiven "1"
    StrCpy $ConfigTouched "1"
  ${EndIf}

  ${GetOptions} $R9 "/FIREWALL=" $R0
  ${IfNot} ${Errors}
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ParseBoolean
    StrCpy $FirewallUI $R1
  ${EndIf}

  ${GetOptions} $R9 "/FIREWALLHTTP=" $R0
  ${IfNot} ${Errors}
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ParseBoolean
    StrCpy $FirewallHTTP $R1
  ${EndIf}

  ${GetOptions} $R9 "/FIREWALLHTTPS=" $R0
  ${IfNot} ${Errors}
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ParseBoolean
    StrCpy $FirewallHTTPS $R1
  ${EndIf}

  ${GetOptions} $R9 "/DESKTOPSHORTCUT=" $R0
  ${IfNot} ${Errors}
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ParseBoolean
    StrCpy $DesktopShortcut $R1
  ${EndIf}

  ${GetOptions} $R9 "/STARTMENUSHORTCUT=" $R0
  ${IfNot} ${Errors}
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ParseBoolean
    StrCpy $StartMenuShortcut $R1
  ${EndIf}

  ${GetOptions} $R9 "/START=" $R0
  ${IfNot} ${Errors}
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ParseBoolean
    StrCpy $ServiceStartNow $R1
  ${EndIf}

  ${GetOptions} $R9 "/OPENUI=" $R0
  ${IfNot} ${Errors}
    StrCpy $R2 $R0
    StrCpy $R0 $R2
    Call ParseBoolean
    StrCpy $OpenUI $R1
  ${EndIf}

  ${GetOptions} $R9 "/CONFIG=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    StrCpy $ConfigImport $R0
  ${EndIf}

  ${GetOptions} $R9 "/INSTALLTYPE=" $R0
  ${IfNot} ${Errors}
  ${AndIf} $R0 != ""
    ${StrCase} $R0 "$R0" "L"
    ${If} $R0 == "express"
      StrCpy $InstallType "express"
    ${Else}
      StrCpy $InstallType "advanced"
    ${EndIf}
  ${EndIf}
FunctionEnd
