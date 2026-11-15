; Windows service, firewall, Event Log and file permission handling.

; Output: $R0 = 1 when the nginx ignition service exists.
Function ServiceExists
  StrCpy $R0 "0"
  ClearErrors
  ReadRegStr $0 HKLM "SYSTEM\CurrentControlSet\Services\${APP_SERVICE}" "ImagePath"
  ${IfNot} ${Errors}
  ${AndIf} $0 != ""
    StrCpy $R0 "1"
  ${EndIf}
FunctionEnd

; Stops the existing service (if any) and records whether it was running.
Function StopService
  StrCpy $PreviousServiceRunning "0"

  Call ServiceExists
  ${If} $R0 != "1"
    Return
  ${EndIf}

  nsExec::ExecToStack `powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$PLUGINSDIR\wait-service.ps1" -Status Running -Check`
  Pop $0
  Pop $1
  ${If} $0 == "0"
    StrCpy $PreviousServiceRunning "1"
  ${EndIf}

  DetailPrint "$(NI_DETAIL_STOP)"

  nsExec::ExecToStack 'sc.exe stop ${APP_SERVICE}'
  Pop $0
  Pop $1

  nsExec::ExecToStack `powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$PLUGINSDIR\wait-service.ps1" -Status Stopped -TimeoutSeconds 60`
  Pop $0
  Pop $1
  ${If} $0 != "0"
  ${AndIf} $0 != "2"
    MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_SERVICE_DELETE)"
    Abort
  ${EndIf}

  ; Also terminate a manually started nginx-ignition.exe; otherwise it can keep
  ; the SQLite database open while the installer rewrites the ACLs.
  nsExec::ExecToStack 'taskkill.exe /F /IM nginx-ignition.exe'
  Pop $0
  Pop $1
FunctionEnd

; Resolves the service account to a name that icacls understands.
Function ResolveServiceAccount
  ${If} $ServiceAccountKind == "custom"
    StrCpy $AccountName "$ServiceAccount"
  ${ElseIf} $ServiceAccountKind == "networkservice"
    StrCpy $AccountName "*S-1-5-20"
  ${ElseIf} $ServiceAccountKind == "localservice"
    StrCpy $AccountName "*S-1-5-19"
  ${Else}
    StrCpy $AccountName "*S-1-5-18"
  ${EndIf}
FunctionEnd

; Grants SYSTEM, Administrators and the service account access to a
; configured path (data, nginx config or VPN config) without changing
; inheritance. Input: $AclPath.
Function GrantAclPath
  ${If} $AclPath == ""
    Return
  ${EndIf}

  nsExec::ExecToStack 'icacls "$AclPath" /grant:r "*S-1-5-18:(OI)(CI)F" /T /C'
  Pop $0
  Pop $1
  nsExec::ExecToStack 'icacls "$AclPath" /grant:r "*S-1-5-32-544:(OI)(CI)F" /T /C'
  Pop $0
  Pop $1

  ${If} $AccountName != "*S-1-5-18"
    nsExec::ExecToStack 'icacls "$AclPath" /grant:r "$AccountName:(OI)(CI)M" /T /C'
    Pop $0
    Pop $1
  ${EndIf}
FunctionEnd

Function ApplyPermissions
  DetailPrint "$(NI_DETAIL_ACL)"

  Call ResolveServiceAccount

  ; Grant each SID with its own icacls call. A single call with several
  ; entries can fail when one of them already exists, and a failed grant must
  ; never be followed by removing inherited access.
  StrCpy $2 "0"

  nsExec::ExecToStack 'icacls "$APPDATA_DIR" /grant:r "*S-1-5-18:(OI)(CI)F" /T /C'
  Pop $0
  Pop $1
  ${If} $0 != "0"
    DetailPrint "$1"
    ; An older build could leave the data tree with an empty or damaged DACL.
    ; Take ownership first, then retry the grant.
    nsExec::ExecToStack 'takeown.exe /F "$APPDATA_DIR" /R /D Y'
    Pop $0
    Pop $1
    nsExec::ExecToStack 'icacls "$APPDATA_DIR" /grant:r "*S-1-5-18:(OI)(CI)F" /T /C'
    Pop $0
    Pop $1
    ${If} $0 != "0"
      DetailPrint "$1"
      StrCpy $2 "1"
    ${EndIf}
  ${EndIf}

  nsExec::ExecToStack 'icacls "$APPDATA_DIR" /grant:r "*S-1-5-32-544:(OI)(CI)F" /T /C'
  Pop $0
  Pop $1
  ${If} $0 != "0"
    DetailPrint "$1"
    StrCpy $2 "1"
  ${EndIf}

  ${If} $AccountName != "*S-1-5-18"
    nsExec::ExecToStack 'icacls "$APPDATA_DIR" /grant:r "$AccountName:(OI)(CI)M" /T /C'
    Pop $0
    Pop $1
    ${If} $0 != "0"
      DetailPrint "$1"
      StrCpy $2 "1"
    ${EndIf}
  ${EndIf}

  ; Only tighten the directory on a fresh install. On an upgrade the ACL is
  ; already in place, and re-applying inheritance removal to an existing
  ; SQLite tree risks leaving the service account without access.
  ${If} $IsUpgrade != "1"
  ${AndIf} $2 == "0"
    nsExec::ExecToStack 'icacls "$APPDATA_DIR" /inheritance:r /T /C'
    Pop $0
    Pop $1
  ${EndIf}

  ; Install directory and configuration file: one SID per call.
  nsExec::ExecToStack 'icacls "$INSTDIR" /grant:r "*S-1-5-18:(OI)(CI)F" /C'
  Pop $0
  Pop $1
  nsExec::ExecToStack 'icacls "$INSTDIR" /grant:r "*S-1-5-32-544:(OI)(CI)F" /C'
  Pop $0
  Pop $1

  ${If} $AccountName != "*S-1-5-18"
    nsExec::ExecToStack 'icacls "$INSTDIR" /grant:r "$AccountName:(OI)(CI)RX" /C'
    Pop $0
    Pop $1
  ${EndIf}

  nsExec::ExecToStack 'icacls "$CONFIG_FILE" /grant:r "*S-1-5-18:F" /C'
  Pop $0
  Pop $1
  nsExec::ExecToStack 'icacls "$CONFIG_FILE" /grant:r "*S-1-5-32-544:F" /C'
  Pop $0
  Pop $1

  ${If} $AccountName != "*S-1-5-18"
    nsExec::ExecToStack 'icacls "$CONFIG_FILE" /grant:r "$AccountName:R" /C'
    Pop $0
    Pop $1

    ${If} $NginxBinaryPath != ""
      nsExec::ExecToStack 'icacls "$NginxBinaryPath" /grant "$AccountName:RX" /C'
      Pop $0
      Pop $1
    ${EndIf}
  ${EndIf}

  ; The data, nginx configuration and VPN directories can be chosen outside
  ; C:\ProgramData, so make sure the service account can use them too.
  StrCpy $AclPath "$DataDir"
  Call GrantAclPath
  StrCpy $AclPath "$NginxConfigPath"
  Call GrantAclPath
  StrCpy $AclPath "$VpnConfigPath"
  Call GrantAclPath
FunctionEnd

Function RegisterEventLog
  DetailPrint "$(NI_DETAIL_EVENTLOG)"
  ; Match the values used by golang.org/x/sys/windows/svc/eventlog.InstallAsEventCreate,
  ; so that the application finds the source already registered and events render
  ; with a valid message file even when the service runs as a low-privileged account.
  WriteRegDWORD HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\${APP_SERVICE}" \
    "CustomSource" 1
  WriteRegExpandStr HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\${APP_SERVICE}" \
    "EventMessageFile" "%SystemRoot%\System32\EventCreate.exe"
  WriteRegDWORD HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\${APP_SERVICE}" \
    "TypesSupported" 7
FunctionEnd

Function ConfigureService
  DetailPrint "$(NI_DETAIL_SERVICE)"

  ${If} $ServiceStartType == "manual"
    StrCpy $StartMode "demand"
  ${Else}
    StrCpy $StartMode "auto"
  ${EndIf}

  Call ServiceExists
  ${If} $R0 == "1"
    ; A normal upgrade keeps the service object so its credentials and current
    ; startup type are preserved; only /SERVICESTART changes it. If a service
    ; exists without the installer's registry entry (fresh install over a
    ; legacy/manual service), the wizard's startup choice must be applied.
    ${If} $IsUpgrade == "1"
    ${AndIf} $ServiceStartArgGiven != "1"
      nsExec::ExecToStack 'sc.exe config ${APP_SERVICE} binPath= "\"$INSTDIR\nginx-ignition.exe\" --config \"$CONFIG_FILE\""'
    ${Else}
      nsExec::ExecToStack 'sc.exe config ${APP_SERVICE} binPath= "\"$INSTDIR\nginx-ignition.exe\" --config \"$CONFIG_FILE\"" start= $StartMode'
    ${EndIf}
    Pop $0
    Pop $1

    ${If} $ServiceAccountArgGiven == "1"
      ${If} $ServiceAccountKind == "custom"
      ${AndIf} $ServiceAccount != ""
      ${AndIf} $ServicePassword != ""
        nsExec::ExecToStack 'sc.exe config ${APP_SERVICE} obj= "$ServiceAccount" password= "$ServicePassword"'
        Pop $0
        Pop $1
        ${If} $0 != "0"
          MessageBox MB_ICONSTOP|MB_OK "$(NI_SERVICE_CREATE_FAILED)$\r$\n$1"
          Abort
        ${EndIf}
      ${ElseIf} $ServiceAccountKind == "networkservice"
        nsExec::ExecToStack 'sc.exe config ${APP_SERVICE} obj= "NT AUTHORITY\NetworkService"'
        Pop $0
        Pop $1
      ${ElseIf} $ServiceAccountKind == "localservice"
        nsExec::ExecToStack 'sc.exe config ${APP_SERVICE} obj= "NT AUTHORITY\LocalService"'
        Pop $0
        Pop $1
      ${ElseIf} $ServiceAccountKind == "localsystem"
        nsExec::ExecToStack 'sc.exe config ${APP_SERVICE} obj= "LocalSystem"'
        Pop $0
        Pop $1
      ${EndIf}
    ${EndIf}

    nsExec::ExecToStack 'sc.exe description ${APP_SERVICE} "$(NI_SERVICE_DESCRIPTION)"'
    Pop $0
    Pop $1
    Return
  ${EndIf}

  ${If} $ServiceAccountKind == "custom"
    nsExec::ExecToStack `powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$PLUGINSDIR\grant-service-logon.ps1" -Account "$ServiceAccount"`
    Pop $0
    Pop $1
    ${If} $0 != "0"
      DetailPrint "$1"
    ${EndIf}

    nsExec::ExecToStack 'sc.exe create ${APP_SERVICE} binPath= "\"$INSTDIR\nginx-ignition.exe\" --config \"$CONFIG_FILE\"" start= $StartMode obj= "$ServiceAccount" password= "$ServicePassword" DisplayName= "${APP_NAME}"'
  ${ElseIf} $ServiceAccountKind == "networkservice"
    nsExec::ExecToStack 'sc.exe create ${APP_SERVICE} binPath= "\"$INSTDIR\nginx-ignition.exe\" --config \"$CONFIG_FILE\"" start= $StartMode obj= "NT AUTHORITY\NetworkService" password= "" DisplayName= "${APP_NAME}"'
  ${ElseIf} $ServiceAccountKind == "localservice"
    nsExec::ExecToStack 'sc.exe create ${APP_SERVICE} binPath= "\"$INSTDIR\nginx-ignition.exe\" --config \"$CONFIG_FILE\"" start= $StartMode obj= "NT AUTHORITY\LocalService" password= "" DisplayName= "${APP_NAME}"'
  ${Else}
    nsExec::ExecToStack 'sc.exe create ${APP_SERVICE} binPath= "\"$INSTDIR\nginx-ignition.exe\" --config \"$CONFIG_FILE\"" start= $StartMode DisplayName= "${APP_NAME}"'
  ${EndIf}

  Pop $0
  Pop $1

  ${If} $0 != "0"
    MessageBox MB_ICONSTOP|MB_OK "$(NI_SERVICE_CREATE_FAILED)$\r$\n$1"
    Abort
  ${EndIf}

  nsExec::ExecToStack 'sc.exe description ${APP_SERVICE} "$(NI_SERVICE_DESCRIPTION)"'
  Pop $0
  Pop $1
FunctionEnd

Function StartService
  ${If} $ServiceStartNow != "1"
  ${AndIf} $PreviousServiceRunning != "1"
    Return
  ${EndIf}

  DetailPrint "$(NI_DETAIL_START)"

  nsExec::ExecToStack 'sc.exe start ${APP_SERVICE}'
  Pop $0
  Pop $1

  nsExec::ExecToStack `powershell.exe -NoProfile -ExecutionPolicy Bypass -File "$PLUGINSDIR\wait-service.ps1" -Status Running -TimeoutSeconds 60`
  Pop $0
  Pop $1

  ${If} $0 != "0"
    MessageBox MB_ICONEXCLAMATION|MB_OK "$(NI_SERVICE_START_FAILED)"
  ${EndIf}
FunctionEnd

Function RemoveFirewallRules
  nsExec::ExecToStack 'netsh advfirewall firewall delete rule name="${APP_RULE_WEB}"'
  Pop $0
  Pop $1
  nsExec::ExecToStack 'netsh advfirewall firewall delete rule name="${APP_RULE_HTTP}"'
  Pop $0
  Pop $1
  nsExec::ExecToStack 'netsh advfirewall firewall delete rule name="${APP_RULE_HTTPS}"'
  Pop $0
  Pop $1
FunctionEnd

Function AddFirewallRule
  ; Input: $R0 = rule name, $R1 = local port
  nsExec::ExecToStack 'netsh advfirewall firewall add rule name="$R0" dir=in action=allow protocol=TCP localport=$R1 profile=domain,private enable=yes'
  Pop $0
  Pop $1
  ${If} $0 != "0"
    StrCpy $FirewallError "1"
  ${EndIf}
FunctionEnd

Function ApplyFirewall
  DetailPrint "$(NI_DETAIL_FIREWALL)"
  StrCpy $FirewallError "0"

  Call RemoveFirewallRules

  ${If} $FirewallUI == "1"
    StrCpy $R0 "${APP_RULE_WEB}"
    StrCpy $R1 "$ServerPort"
    Call AddFirewallRule
  ${EndIf}

  ${If} $FirewallHTTP == "1"
    StrCpy $R0 "${APP_RULE_HTTP}"
    StrCpy $R1 "80"
    Call AddFirewallRule
  ${EndIf}

  ${If} $FirewallHTTPS == "1"
    StrCpy $R0 "${APP_RULE_HTTPS}"
    StrCpy $R1 "443"
    Call AddFirewallRule
  ${EndIf}

  ${If} $FirewallError == "1"
    MessageBox MB_ICONEXCLAMATION|MB_OK "$(NI_ERR_FIREWALL)"
  ${EndIf}
FunctionEnd
