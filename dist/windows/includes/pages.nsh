; Custom wizard pages for the nginx ignition installer.

Var hInstallTypeExpress
Var hInstallTypeAdvanced

Var hNginxBinary
Var hNginxConfig
Var hVpnConfig

Var hServerAddress
Var hServerPort
Var hDataDir

Var hDbSqlite
Var hDbPostgres
Var hDbHost
Var hDbPort
Var hDbName
Var hDbSchema
Var hDbUsername
Var hDbPassword
Var hDbSslMode

Var hSvcLocalSystem
Var hSvcNetworkService
Var hSvcLocalService
Var hSvcCustom
Var hSvcStartAuto
Var hSvcStartManual
Var hSvcStartNow

Var hCustomUser
Var hCustomPassword
Var hCustomPasswordConfirm

Var hFwUI
Var hFwHTTP
Var hFwHTTPS
Var hDesktopShortcut
Var hStartMenuShortcut

;--------------------------------
;Installation type

Function InstallTypePageCreate
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}

  !insertmacro MUI_HEADER_TEXT "$(NI_INSTALLTYPE_TITLE)" "$(NI_INSTALLTYPE_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateRadioButton} 0 0 100% 12u "$(NI_INSTALLTYPE_EXPRESS)"
  Pop $hInstallTypeExpress
  ${NSD_AddStyle} $hInstallTypeExpress ${WS_GROUP}

  ${NSD_CreateLabel} 12u 12u 268u 32u "$(NI_INSTALLTYPE_EXPRESS_DESC)"
  Pop $0

  ${NSD_CreateRadioButton} 0 50u 100% 12u "$(NI_INSTALLTYPE_ADVANCED)"
  Pop $hInstallTypeAdvanced

  ${NSD_CreateLabel} 12u 62u 268u 28u "$(NI_INSTALLTYPE_ADVANCED_DESC)"
  Pop $0

  ${If} $InstallType == "advanced"
    ${NSD_Check} $hInstallTypeAdvanced
  ${Else}
    ${NSD_Check} $hInstallTypeExpress
  ${EndIf}

  nsDialogs::Show
FunctionEnd

Function InstallTypePageLeave
  ${NSD_GetState} $hInstallTypeAdvanced $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $InstallType "advanced"
    Call DetectNginx
  ${Else}
    StrCpy $InstallType "express"
    Call ApplyExpressDefaults
  ${EndIf}
FunctionEnd

;--------------------------------
;nginx configuration

Function NginxPageCreate
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}
  ${If} $InstallType != "advanced"
    Abort
  ${EndIf}

  !insertmacro MUI_HEADER_TEXT "$(NI_NGINX_TITLE)" "$(NI_NGINX_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateGroupBox} 0 0 300u 86u "$(NI_GROUP_NGINX)"
  Pop $0

  ${NSD_CreateLabel} 10u 12u 280u 9u "$(NI_NGINX_PATH)"
  Pop $0
  ${NSD_CreateText} 10u 21u 210u 12u "$NginxBinaryPath"
  Pop $hNginxBinary
  ${NSD_CreateButton} 225u 21u 65u 12u "$(NI_BROWSE)"
  Pop $0
  ${NSD_OnClick} $0 BrowseNginxBinary

  ${NSD_CreateLabel} 10u 37u 280u 9u "$(NI_NGINX_CONFIG_DIR)"
  Pop $0
  ${NSD_CreateText} 10u 46u 210u 12u "$NginxConfigPath"
  Pop $hNginxConfig
  ${NSD_CreateButton} 225u 46u 65u 12u "$(NI_BROWSE)"
  Pop $0
  ${NSD_OnClick} $0 BrowseNginxConfig

  ${NSD_CreateLabel} 10u 62u 280u 20u "$(NI_NGINX_HINT)"
  Pop $0

  ${NSD_CreateGroupBox} 0 92u 300u 36u "$(NI_GROUP_VPN)"
  Pop $0

  ${NSD_CreateLabel} 10u 103u 280u 9u "$(NI_VPN_CONFIG_DIR)"
  Pop $0
  ${NSD_CreateText} 10u 112u 210u 12u "$VpnConfigPath"
  Pop $hVpnConfig
  ${NSD_CreateButton} 225u 112u 65u 12u "$(NI_BROWSE)"
  Pop $0
  ${NSD_OnClick} $0 BrowseVpnConfig

  nsDialogs::Show
FunctionEnd

Function NginxPageLeave
  ${NSD_GetText} $hNginxBinary $NginxBinaryPath
  ${NSD_GetText} $hNginxConfig $NginxConfigPath
  ${NSD_GetText} $hVpnConfig $VpnConfigPath

  ${If} $NginxConfigPath == ""
    StrCpy $NginxConfigPath "$APPDATA_DIR\nginx"
  ${EndIf}
  ${If} $VpnConfigPath == ""
    StrCpy $VpnConfigPath "$APPDATA_DIR\vpn"
  ${EndIf}

  ${If} $NginxBinaryPath == ""
    MessageBox MB_YESNO|MB_ICONEXCLAMATION "$(NI_NGINX_NOT_FOUND)" /SD IDYES IDNO stop
    Return
  ${EndIf}

  ${IfNot} ${FileExists} "$NginxBinaryPath"
    MessageBox MB_YESNO|MB_ICONEXCLAMATION "$(NI_NGINX_NOT_FOUND)" /SD IDYES IDNO stop
  ${EndIf}

  Return

  stop:
    Abort
FunctionEnd

;--------------------------------
;Server and data

Function ServerPageCreate
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}
  ${If} $InstallType != "advanced"
    Abort
  ${EndIf}

  !insertmacro MUI_HEADER_TEXT "$(NI_SERVER_TITLE)" "$(NI_SERVER_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateGroupBox} 0 0 300u 66u "$(NI_GROUP_WEB)"
  Pop $0

  ${NSD_CreateLabel} 10u 14u 120u 9u "$(NI_SERVER_ADDRESS)"
  Pop $0
  ${NSD_CreateText} 10u 24u 180u 12u "$ServerAddress"
  Pop $hServerAddress

  ${NSD_CreateLabel} 200u 14u 40u 9u "$(NI_SERVER_PORT)"
  Pop $0
  ${NSD_CreateNumber} 200u 24u 90u 12u "$ServerPort"
  Pop $hServerPort

  ${NSD_CreateLabel} 10u 42u 280u 20u "$(NI_SERVER_HINT)"
  Pop $0

  ${NSD_CreateGroupBox} 0 74u 300u 44u "$(NI_GROUP_DATA)"
  Pop $0

  ${NSD_CreateLabel} 10u 87u 280u 9u "$(NI_DATA_DIR)"
  Pop $0
  ${NSD_CreateText} 10u 97u 210u 12u "$DataDir"
  Pop $hDataDir
  ${NSD_CreateButton} 225u 97u 65u 12u "$(NI_BROWSE)"
  Pop $0
  ${NSD_OnClick} $0 BrowseDataDir

  nsDialogs::Show
FunctionEnd

Function ServerPageLeave
  ${NSD_GetText} $hServerAddress $ServerAddress
  ${NSD_GetText} $hServerPort $ServerPort
  ${NSD_GetText} $hDataDir $DataDir

  ${If} $ServerAddress == ""
    StrCpy $ServerAddress "127.0.0.1"
  ${EndIf}
  ${If} $DataDir == ""
    StrCpy $DataDir "$APPDATA_DIR\data"
  ${EndIf}

  StrCpy $R0 "$ServerPort"
  Call ValidatePort
  ${If} $R1 != "1"
    MessageBox MB_ICONEXCLAMATION|MB_OK "$(NI_INVALID_PORT)"
    Abort
  ${EndIf}
FunctionEnd

;--------------------------------
;Database

Function DatabasePageCreate
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}
  ${If} $InstallType != "advanced"
    Abort
  ${EndIf}

  !insertmacro MUI_HEADER_TEXT "$(NI_DATABASE_TITLE)" "$(NI_DATABASE_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateGroupBox} 0 0 300u 42u "$(NI_DATABASE_TITLE)"
  Pop $0

  ${NSD_CreateRadioButton} 10u 14u 280u 11u "$(NI_DB_SQLITE)"
  Pop $hDbSqlite
  ${NSD_AddStyle} $hDbSqlite ${WS_GROUP}
  ${NSD_CreateRadioButton} 10u 27u 280u 11u "$(NI_DB_POSTGRES)"
  Pop $hDbPostgres

  ${NSD_CreateGroupBox} 0 48u 300u 72u "$(NI_GROUP_POSTGRES)"
  Pop $0

  ${NSD_CreateLabel} 10u 62u 40u 9u "$(NI_DB_HOST)"
  Pop $0
  ${NSD_CreateText} 55u 61u 90u 12u "$DbHost"
  Pop $hDbHost

  ${NSD_CreateLabel} 155u 62u 28u 9u "$(NI_DB_PORT)"
  Pop $0
  ${NSD_CreateNumber} 188u 61u 90u 12u "$DbPort"
  Pop $hDbPort

  ${NSD_CreateLabel} 10u 76u 40u 9u "$(NI_DB_NAME)"
  Pop $0
  ${NSD_CreateText} 55u 75u 90u 12u "$DbName"
  Pop $hDbName

  ${NSD_CreateLabel} 155u 76u 28u 9u "$(NI_DB_SCHEMA)"
  Pop $0
  ${NSD_CreateText} 188u 75u 90u 12u "$DbSchema"
  Pop $hDbSchema

  ${NSD_CreateLabel} 10u 90u 40u 9u "$(NI_DB_USER)"
  Pop $0
  ${NSD_CreateText} 55u 89u 90u 12u "$DbUsername"
  Pop $hDbUsername

  ${NSD_CreateLabel} 10u 104u 40u 9u "$(NI_DB_PASSWORD)"
  Pop $0
  ${NSD_CreatePassword} 55u 103u 90u 12u "$DbPassword"
  Pop $hDbPassword

  ${NSD_CreateLabel} 155u 104u 40u 9u "$(NI_DB_SSLMODE)"
  Pop $0
  ${NSD_CreateDropList} 200u 103u 90u 50u ""
  Pop $hDbSslMode
  ${NSD_CB_AddString} $hDbSslMode "$(NI_DB_SSL_DISABLE)"
  ${NSD_CB_AddString} $hDbSslMode "$(NI_DB_SSL_REQUIRE)"
  ${If} $DbSslMode == "disable"
    ${NSD_CB_SelectString} $hDbSslMode "$(NI_DB_SSL_DISABLE)"
  ${Else}
    ${NSD_CB_SelectString} $hDbSslMode "$(NI_DB_SSL_REQUIRE)"
  ${EndIf}

  ${If} $DatabaseDriver == "postgres"
    ${NSD_Check} $hDbPostgres
  ${Else}
    ${NSD_Check} $hDbSqlite
  ${EndIf}

  ${NSD_OnClick} $hDbSqlite DatabaseUpdateControls
  ${NSD_OnClick} $hDbPostgres DatabaseUpdateControls
  Call DatabaseUpdateControls

  nsDialogs::Show
FunctionEnd

Function DatabaseUpdateControls
  ${NSD_GetState} $hDbPostgres $0
  ${If} $0 == ${BST_CHECKED}
    EnableWindow $hDbHost 1
    EnableWindow $hDbPort 1
    EnableWindow $hDbName 1
    EnableWindow $hDbSchema 1
    EnableWindow $hDbUsername 1
    EnableWindow $hDbPassword 1
    EnableWindow $hDbSslMode 1
  ${Else}
    EnableWindow $hDbHost 0
    EnableWindow $hDbPort 0
    EnableWindow $hDbName 0
    EnableWindow $hDbSchema 0
    EnableWindow $hDbUsername 0
    EnableWindow $hDbPassword 0
    EnableWindow $hDbSslMode 0
  ${EndIf}
FunctionEnd

Function DatabasePageLeave
  ${NSD_GetState} $hDbPostgres $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $DatabaseDriver "postgres"

    ${NSD_GetText} $hDbHost $DbHost
    ${NSD_GetText} $hDbPort $DbPort
    ${NSD_GetText} $hDbName $DbName
    ${NSD_GetText} $hDbSchema $DbSchema
    ${NSD_GetText} $hDbUsername $DbUsername
    ${NSD_GetText} $hDbPassword $0
    ${If} $0 != ""
      StrCpy $DbPassword "$0"
    ${ElseIf} $ExistingDbPassword != ""
      StrCpy $DbPassword "$ExistingDbPassword"
    ${EndIf}

    SendMessage $hDbSslMode ${CB_GETCURSEL} 0 0
    Pop $0
    ${If} $0 == 0
      StrCpy $DbSslMode "disable"
    ${Else}
      StrCpy $DbSslMode "require"
    ${EndIf}

    ${If} $DbHost == ""
    ${OrIf} $DbName == ""
    ${OrIf} $DbUsername == ""
    ${OrIf} $DbPassword == ""
      MessageBox MB_ICONEXCLAMATION|MB_OK "$(NI_DB_REQUIRED)"
      Abort
    ${EndIf}

    StrCpy $R0 "$DbPort"
    Call ValidatePort
    ${If} $R1 != "1"
      MessageBox MB_ICONEXCLAMATION|MB_OK "$(NI_INVALID_PORT)"
      Abort
    ${EndIf}

    ${If} $DbSchema == ""
      StrCpy $DbSchema "public"
    ${EndIf}

    ${If} $IsUpgrade == "1"
    ${AndIf} $PreviousDriver != ""
    ${AndIf} $PreviousDriver != "postgres"
      MessageBox MB_YESNO|MB_ICONQUESTION "$(NI_DB_CHANGE_WARNING)" /SD IDNO IDNO stop
    ${EndIf}
    Return

    stop:
      Abort
  ${Else}
    StrCpy $DatabaseDriver "sqlite"
  ${EndIf}
FunctionEnd

;--------------------------------
;Windows service

Function ServicePageCreate
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}
  ${If} $InstallType != "advanced"
    Abort
  ${EndIf}

  !insertmacro MUI_HEADER_TEXT "$(NI_SERVICE_TITLE)" "$(NI_SERVICE_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateGroupBox} 0 0 300u 66u "$(NI_SERVICE_ACCOUNT)"
  Pop $0

  ${NSD_CreateRadioButton} 10u 14u 280u 10u "$(NI_SERVICE_LOCALSYSTEM)"
  Pop $hSvcLocalSystem
  ${NSD_AddStyle} $hSvcLocalSystem ${WS_GROUP}
  ${NSD_CreateRadioButton} 10u 26u 280u 10u "$(NI_SERVICE_NETWORKSERVICE)"
  Pop $hSvcNetworkService
  ${NSD_CreateRadioButton} 10u 38u 280u 10u "$(NI_SERVICE_LOCALSERVICE)"
  Pop $hSvcLocalService
  ${NSD_CreateRadioButton} 10u 50u 280u 10u "$(NI_SERVICE_CUSTOM)"
  Pop $hSvcCustom

  ${NSD_CreateGroupBox} 0 72u 300u 46u "$(NI_SERVICE_STARTTYPE)"
  Pop $0

  ${NSD_CreateRadioButton} 10u 86u 50u 10u "$(NI_SERVICE_AUTO)"
  Pop $hSvcStartAuto
  ${NSD_AddStyle} $hSvcStartAuto ${WS_GROUP}
  ${NSD_CreateRadioButton} 65u 86u 60u 10u "$(NI_SERVICE_MANUAL)"
  Pop $hSvcStartManual

  ${NSD_CreateCheckbox} 10u 102u 280u 10u "$(NI_SERVICE_START_NOW)"
  Pop $hSvcStartNow

  ${If} $ServiceAccountKind == "networkservice"
    ${NSD_Check} $hSvcNetworkService
  ${ElseIf} $ServiceAccountKind == "localservice"
    ${NSD_Check} $hSvcLocalService
  ${ElseIf} $ServiceAccountKind == "custom"
    ${NSD_Check} $hSvcCustom
  ${Else}
    ${NSD_Check} $hSvcLocalSystem
  ${EndIf}

  ${If} $ServiceStartType == "manual"
    ${NSD_Check} $hSvcStartManual
  ${Else}
    ${NSD_Check} $hSvcStartAuto
  ${EndIf}

  ${If} $ServiceStartNow == "1"
    ${NSD_Check} $hSvcStartNow
  ${EndIf}

  nsDialogs::Show
FunctionEnd

Function ServicePageLeave
  ${NSD_GetState} $hSvcCustom $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $ServiceAccountKind "custom"
  ${Else}
    StrCpy $ServiceAccount ""
    StrCpy $ServicePassword ""

    ${NSD_GetState} $hSvcNetworkService $0
    ${If} $0 == ${BST_CHECKED}
      StrCpy $ServiceAccountKind "networkservice"
    ${Else}
      ${NSD_GetState} $hSvcLocalService $0
      ${If} $0 == ${BST_CHECKED}
        StrCpy $ServiceAccountKind "localservice"
      ${Else}
        StrCpy $ServiceAccountKind "localsystem"
      ${EndIf}
    ${EndIf}
  ${EndIf}

  ${NSD_GetState} $hSvcStartManual $0
  ${If} $0 == ${BST_CHECKED}
    StrCpy $ServiceStartType "manual"
  ${Else}
    StrCpy $ServiceStartType "auto"
  ${EndIf}

  ${NSD_GetState} $hSvcStartNow $ServiceStartNow
FunctionEnd

;--------------------------------
;Custom service account

Function CustomAccountPageCreate
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}
  ${If} $InstallType != "advanced"
    Abort
  ${EndIf}
  ${If} $ServiceAccountKind != "custom"
    Abort
  ${EndIf}

  !insertmacro MUI_HEADER_TEXT "$(NI_SERVICE_CUSTOM)" "$(NI_SERVICE_CUSTOM_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateGroupBox} 0 0 300u 62u "$(NI_SERVICE_CUSTOM)"
  Pop $0

  ${NSD_CreateLabel} 10u 14u 60u 9u "$(NI_SERVICE_USER)"
  Pop $0
  ${NSD_CreateText} 75u 13u 215u 12u "$ServiceAccount"
  Pop $hCustomUser

  ${NSD_CreateLabel} 10u 30u 60u 9u "$(NI_SERVICE_PASSWORD)"
  Pop $0
  ${NSD_CreatePassword} 75u 29u 215u 12u "$ServicePassword"
  Pop $hCustomPassword

  ${NSD_CreateLabel} 10u 46u 60u 9u "$(NI_SERVICE_PASSWORD_CONFIRM)"
  Pop $0
  ${NSD_CreatePassword} 75u 45u 215u 12u "$ServicePassword"
  Pop $hCustomPasswordConfirm

  ${NSD_CreateLabel} 0 70u 300u 26u "$(NI_SERVICE_CUSTOM_HINT)"
  Pop $0

  nsDialogs::Show
FunctionEnd

Function CustomAccountPageLeave
  ${NSD_GetText} $hCustomUser $ServiceAccount
  ${NSD_GetText} $hCustomPassword $ServicePassword
  ${NSD_GetText} $hCustomPasswordConfirm $0

  ${If} $ServiceAccount == ""
  ${OrIf} $ServicePassword == ""
    MessageBox MB_ICONEXCLAMATION|MB_OK "$(NI_SERVICE_CUSTOM_REQUIRED)"
    Abort
  ${EndIf}

  ${If} $ServicePassword != $0
    MessageBox MB_ICONEXCLAMATION|MB_OK "$(NI_SERVICE_PASSWORD_MISMATCH)"
    Abort
  ${EndIf}
FunctionEnd

;--------------------------------
;Path pickers

Function BrowseNginxBinary
  Pop $0
  ${NSD_GetText} $hNginxBinary $0
  nsDialogs::SelectFileDialog open "$0" "$(NI_FILE_FILTER)"
  Pop $0
  ${If} $0 != ""
    ${NSD_SetText} $hNginxBinary "$0"
  ${EndIf}
FunctionEnd

Function BrowseNginxConfig
  Pop $0
  ${NSD_GetText} $hNginxConfig $0
  nsDialogs::SelectFolderDialog "$(NI_GROUP_NGINX)" "$0"
  Pop $0
  ${If} $0 != ""
  ${AndIf} $0 != "error"
    ${NSD_SetText} $hNginxConfig "$0"
  ${EndIf}
FunctionEnd

Function BrowseVpnConfig
  Pop $0
  ${NSD_GetText} $hVpnConfig $0
  nsDialogs::SelectFolderDialog "$(NI_GROUP_VPN)" "$0"
  Pop $0
  ${If} $0 != ""
  ${AndIf} $0 != "error"
    ${NSD_SetText} $hVpnConfig "$0"
  ${EndIf}
FunctionEnd

Function BrowseDataDir
  Pop $0
  ${NSD_GetText} $hDataDir $0
  nsDialogs::SelectFolderDialog "$(NI_GROUP_DATA)" "$0"
  Pop $0
  ${If} $0 != ""
  ${AndIf} $0 != "error"
    ${NSD_SetText} $hDataDir "$0"
  ${EndIf}
FunctionEnd

;--------------------------------
;Firewall and shortcuts

Function FirewallPageCreate
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}
  ${If} $InstallType != "advanced"
    Abort
  ${EndIf}

  !insertmacro MUI_HEADER_TEXT "$(NI_FIREWALL_TITLE)" "$(NI_FIREWALL_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateGroupBox} 0 0 300u 56u "$(NI_GROUP_FIREWALL)"
  Pop $0

  ${NSD_CreateCheckbox} 10u 14u 280u 10u "$(NI_FIREWALL_UI)"
  Pop $hFwUI
  ${If} $FirewallUI == "1"
    ${NSD_Check} $hFwUI
  ${EndIf}

  ${NSD_CreateCheckbox} 10u 28u 280u 10u "$(NI_FIREWALL_HTTP)"
  Pop $hFwHTTP
  ${If} $FirewallHTTP == "1"
    ${NSD_Check} $hFwHTTP
  ${EndIf}

  ${NSD_CreateCheckbox} 10u 42u 280u 10u "$(NI_FIREWALL_HTTPS)"
  Pop $hFwHTTPS
  ${If} $FirewallHTTPS == "1"
    ${NSD_Check} $hFwHTTPS
  ${EndIf}

  ${NSD_CreateGroupBox} 0 62u 300u 44u "$(NI_GROUP_SHORTCUTS)"
  Pop $0

  ${NSD_CreateCheckbox} 10u 76u 280u 10u "$(NI_DESKTOP_SHORTCUT)"
  Pop $hDesktopShortcut
  ${If} $DesktopShortcut == "1"
    ${NSD_Check} $hDesktopShortcut
  ${EndIf}

  ${NSD_CreateCheckbox} 10u 90u 280u 10u "$(NI_STARTMENU_SHORTCUT)"
  Pop $hStartMenuShortcut
  ${If} $StartMenuShortcut == "1"
    ${NSD_Check} $hStartMenuShortcut
  ${EndIf}

  nsDialogs::Show
FunctionEnd

Function FirewallPageLeave
  ${NSD_GetState} $hFwUI $FirewallUI
  ${NSD_GetState} $hFwHTTP $FirewallHTTP
  ${NSD_GetState} $hFwHTTPS $FirewallHTTPS
  ${NSD_GetState} $hDesktopShortcut $DesktopShortcut
  ${NSD_GetState} $hStartMenuShortcut $StartMenuShortcut
FunctionEnd
