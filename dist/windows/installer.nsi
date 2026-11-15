; nginx ignition Windows installer.
;
; Build with:
;   makensis -DARCH=amd64 -DVERSION=2.51.0 installer.nsi
;
; Optional defines:
;   BINARY       path to the nginx-ignition.exe to package
;   OUTFILE      output installer path
;   ASSETSDIR    directory with the ICO and BMP artwork
;   LANGDIR      directory with the language files
;   SCRIPTDIR    directory with the PowerShell helper scripts
;   TEMPLATEDIR  directory with the properties templates
;   ROOT         repository root (used for the LICENSE file)
;   NSIS_COMPRESS compressor to use, defaults to lzma

!ifndef ARCH
  !error "ARCH must be defined (amd64 or arm64)"
!endif
!ifndef VERSION
  !error "VERSION must be defined"
!endif

!ifndef ASSETSDIR
  !define ASSETSDIR "assets"
!endif
!ifndef LANGDIR
  !define LANGDIR "lang"
!endif
!ifndef SCRIPTDIR
  !define SCRIPTDIR "scripts"
!endif
!ifndef TEMPLATEDIR
  !define TEMPLATEDIR "templates"
!endif
!ifndef ROOT
  !define ROOT "../.."
!endif
!ifndef BINARY
  !define BINARY "${ROOT}/build/windows/${ARCH}.exe"
!endif
!ifndef OUTFILE
  !define OUTFILE "${ROOT}/build/nginx-ignition-${VERSION}.windows-${ARCH}-setup.exe"
!endif
!ifndef NSIS_COMPRESS
  !define NSIS_COMPRESS "lzma"
!endif

!define APP_NAME "nginx ignition"
!define APP_ID "nginx-ignition"
!define APP_REGKEY "Software\${APP_ID}"
!define APP_UNINSTALLKEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_ID}"
!define APP_SERVICE "nginx-ignition"
!define APP_URL "https://nginxignition.app"
!define APP_RULE_WEB "nginx ignition (web interface)"
!define APP_RULE_HTTP "nginx ignition (HTTP)"
!define APP_RULE_HTTPS "nginx ignition (HTTPS)"

Unicode true
ManifestDPIAware true
ManifestSupportedOS all
RequestExecutionLevel admin

Name "$(NI_NAME)"
Caption "${APP_NAME} ${VERSION}"
BrandingText "${APP_NAME} ${VERSION}"
OutFile "${OUTFILE}"
InstallDir "$PROGRAMFILES64\${APP_ID}"
InstallDirRegKey HKLM "${APP_REGKEY}" "InstallDir"
ShowInstDetails show
ShowUninstDetails show
CRCCheck on

SetCompressor ${NSIS_COMPRESS}
!if "${NSIS_COMPRESS}" != "zlib"
  SetCompressorDictSize 64
!endif

!include "MUI2.nsh"
!include "LogicLib.nsh"
!include "FileFunc.nsh"
!include "WinVer.nsh"
!include "x64.nsh"
!include "nsDialogs.nsh"
!include "StrFunc.nsh"
${Using:StrFunc} StrTrimNewLines
${Using:StrFunc} StrLoc
${Using:StrFunc} StrStr
${Using:StrFunc} StrCase

!define MUI_ICON "${ASSETSDIR}/nginx-ignition.ico"
!define MUI_UNICON "${ASSETSDIR}/nginx-ignition.ico"
!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_RIGHT
; Use the 4x artwork so NSIS can downscale it with LoadImage/RESIZETOFIT.
; This keeps the wizard sharp on high-DPI displays (up to 400% scaling)
; instead of stretching a 1x bitmap.
!define MUI_HEADERIMAGE_BITMAP "${ASSETSDIR}/header@4x.bmp"
!define MUI_WELCOMEFINISHPAGE_BITMAP "${ASSETSDIR}/sidebar@4x.bmp"
!define MUI_UNWELCOMEFINISHPAGE_BITMAP "${ASSETSDIR}/sidebar@4x.bmp"
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_LINK_LOCATION "${APP_URL}"
!define MUI_FINISHPAGE_LINK "$(NI_FINISH_LINK)"
!define MUI_FINISHPAGE_RUN ""
!define MUI_FINISHPAGE_RUN_FUNCTION "OpenUI"
!define MUI_FINISHPAGE_RUN_TEXT "$(NI_FINISH_RUN)"
!define MUI_FINISHPAGE_TITLE "$(NI_FINISH_TITLE)"
!define MUI_FINISHPAGE_TEXT "$FinishText"
!define MUI_FINISHPAGE_TEXT_LARGE
!define MUI_LANGDLL_REGISTRY_ROOT "HKLM"
!define MUI_LANGDLL_REGISTRY_KEY "${APP_REGKEY}"
!define MUI_LANGDLL_REGISTRY_VALUENAME "InstallerLanguage"

VIProductVersion "${VERSION}.0"
VIAddVersionKey "ProductName" "${APP_NAME}"
VIAddVersionKey "ProductVersion" "${VERSION}"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "FileDescription" "${APP_NAME} setup"
VIAddVersionKey "LegalCopyright" "Copyright © Lucas Dillmann"
VIAddVersionKey "CompanyName" "Lucas Dillmann"

Var APPDATA_DIR
Var CONFIG_FILE
Var InstallType
Var NginxBinaryPath
Var NginxConfigPath
Var VpnConfigPath
Var ServerAddress
Var ServerPort
Var DataDir
Var DatabaseDriver
Var DbHost
Var DbPort
Var DbName
Var DbSchema
Var DbUsername
Var DbPassword
Var DbSslMode
Var ExistingDbPassword
Var PreviousDriver
Var ServiceAccountKind
Var ServiceAccount
Var ServicePassword
Var ServiceStartType
Var ServiceStartNow
Var ServiceAccountArgGiven
Var ServiceStartArgGiven
Var PreviousServiceRunning
Var FirewallUI
Var FirewallHTTP
Var FirewallHTTPS
Var FirewallError
Var DesktopShortcut
Var StartMenuShortcut
Var OpenUI
Var ConfigImport
Var ConfigTouched
Var IsUpgrade
Var HasExistingConfig
Var ConfigLoadPath
Var AclPath
Var WelcomeText
Var FinishText
Var UiUrl
Var StartMode
Var AccountName
Var UninstallPurge
Var hUninstallPurge

!include "includes\checks.nsh"
!include "includes\config.nsh"
!include "includes\service.nsh"
!include "includes\pages.nsh"

;--------------------------------
;Installer pages

!define MUI_WELCOMEPAGE_TITLE "$(NI_WELCOME_TITLE)"
!define MUI_WELCOMEPAGE_TEXT "$WelcomeText"
!define MUI_PAGE_CUSTOMFUNCTION_PRE WelcomePagePre
!insertmacro MUI_PAGE_WELCOME

!define MUI_PAGE_CUSTOMFUNCTION_PRE SkipOnUpgrade
!insertmacro MUI_PAGE_LICENSE "${ROOT}/LICENSE"

Page custom InstallTypePageCreate InstallTypePageLeave

!define MUI_PAGE_CUSTOMFUNCTION_PRE SkipOnUpgrade
!insertmacro MUI_PAGE_DIRECTORY

Page custom NginxPageCreate NginxPageLeave
Page custom ServerPageCreate ServerPageLeave
Page custom DatabasePageCreate DatabasePageLeave
Page custom ServicePageCreate ServicePageLeave
Page custom CustomAccountPageCreate CustomAccountPageLeave
Page custom FirewallPageCreate FirewallPageLeave

!insertmacro MUI_PAGE_INSTFILES

!define MUI_PAGE_CUSTOMFUNCTION_PRE FinishPagePre
!insertmacro MUI_PAGE_FINISH

;--------------------------------
;Uninstaller pages

!insertmacro MUI_UNPAGE_CONFIRM
UninstPage custom un.OptionsPageCreate un.OptionsPageLeave
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

;--------------------------------
;Languages

!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "PortugueseBR"
!insertmacro MUI_LANGUAGE "SimpChinese"
!insertmacro MUI_LANGUAGE "German"
!insertmacro MUI_LANGUAGE "French"
!insertmacro MUI_LANGUAGE "SpanishInternational"
!insertmacro MUI_LANGUAGE "Vietnamese"
!insertmacro MUI_LANGUAGE "Russian"
!insertmacro MUI_LANGUAGE "Japanese"

!insertmacro MUI_RESERVEFILE_LANGDLL

!include "${LANGDIR}\en.nsh"
!include "${LANGDIR}\pt.nsh"
!include "${LANGDIR}\zh.nsh"
!include "${LANGDIR}\de.nsh"
!include "${LANGDIR}\fr.nsh"
!include "${LANGDIR}\es.nsh"
!include "${LANGDIR}\vi.nsh"
!include "${LANGDIR}\ru.nsh"
!include "${LANGDIR}\ja.nsh"

;--------------------------------
;Initialization

Function .onInit
  !insertmacro MUI_LANGDLL_DISPLAY

  SetShellVarContext all

  StrCpy $APPDATA_DIR "$APPDATA\${APP_ID}"
  StrCpy $CONFIG_FILE "$APPDATA_DIR\nginx-ignition.properties"

  StrCpy $InstallType "express"
  StrCpy $NginxBinaryPath ""
  StrCpy $NginxConfigPath "$APPDATA_DIR\nginx"
  StrCpy $VpnConfigPath "$APPDATA_DIR\vpn"
  StrCpy $ServerAddress "127.0.0.1"
  StrCpy $ServerPort "8090"
  StrCpy $DataDir "$APPDATA_DIR\data"
  StrCpy $DatabaseDriver "sqlite"
  StrCpy $DbHost ""
  StrCpy $DbPort "5432"
  StrCpy $DbName ""
  StrCpy $DbSchema "public"
  StrCpy $DbUsername ""
  StrCpy $DbPassword ""
  StrCpy $DbSslMode "require"
  StrCpy $ExistingDbPassword ""
  StrCpy $PreviousDriver ""
  StrCpy $ServiceAccountKind "localsystem"
  StrCpy $ServiceAccount ""
  StrCpy $ServicePassword ""
  StrCpy $ServiceStartType "auto"
  StrCpy $ServiceStartNow "1"
  StrCpy $ServiceAccountArgGiven "0"
  StrCpy $ServiceStartArgGiven "0"
  StrCpy $PreviousServiceRunning "0"
  StrCpy $FirewallUI "1"
  StrCpy $FirewallHTTP "0"
  StrCpy $FirewallHTTPS "0"
  StrCpy $FirewallError "0"
  StrCpy $DesktopShortcut "0"
  StrCpy $StartMenuShortcut "1"
  StrCpy $OpenUI "1"
  StrCpy $ConfigImport ""
  StrCpy $ConfigTouched "0"
  StrCpy $IsUpgrade "0"
  StrCpy $HasExistingConfig "0"
  StrCpy $ConfigLoadPath ""

  InitPluginsDir
  File "/oname=$PLUGINSDIR\update-config.ps1" "${SCRIPTDIR}/update-config.ps1"
  File "/oname=$PLUGINSDIR\read-config.ps1" "${SCRIPTDIR}/read-config.ps1"
  File "/oname=$PLUGINSDIR\wait-service.ps1" "${SCRIPTDIR}/wait-service.ps1"
  File "/oname=$PLUGINSDIR\grant-service-logon.ps1" "${SCRIPTDIR}/grant-service-logon.ps1"
  File "/oname=$PLUGINSDIR\nginx-ignition.sqlite.properties" "${TEMPLATEDIR}/nginx-ignition.sqlite.properties.template"
  File "/oname=$PLUGINSDIR\nginx-ignition.postgres.properties" "${TEMPLATEDIR}/nginx-ignition.postgres.properties.template"

  Call CheckPrerequisites

  ReadRegStr $0 HKLM "${APP_REGKEY}" "InstallDir"
  ${If} $0 != ""
    StrCpy $IsUpgrade "1"
    StrCpy $INSTDIR $0
    StrCpy $InstallType "advanced"
    StrCpy $ServiceStartNow "0"

    ReadRegStr $0 HKLM "${APP_REGKEY}" "ConfigFile"
    ${If} $0 != ""
      StrCpy $CONFIG_FILE $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "DataDir"
    ${If} $0 != ""
      StrCpy $DataDir $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "ServiceAccountKind"
    ${If} $0 != ""
      StrCpy $ServiceAccountKind $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "ServiceAccount"
    ${If} $0 != ""
      StrCpy $ServiceAccount $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "ServiceStartType"
    ${If} $0 != ""
      StrCpy $ServiceStartType $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "FirewallUI"
    ${If} $0 != ""
      StrCpy $FirewallUI $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "FirewallHTTP"
    ${If} $0 != ""
      StrCpy $FirewallHTTP $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "FirewallHTTPS"
    ${If} $0 != ""
      StrCpy $FirewallHTTPS $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "DesktopShortcut"
    ${If} $0 != ""
      StrCpy $DesktopShortcut $0
    ${EndIf}

    ReadRegStr $0 HKLM "${APP_REGKEY}" "StartMenuShortcut"
    ${If} $0 != ""
      StrCpy $StartMenuShortcut $0
    ${EndIf}

    StrCpy $HasExistingConfig "1"
    StrCpy $ConfigLoadPath "$CONFIG_FILE"
    Call LoadExistingConfig
  ${Else}
    ${If} ${FileExists} "$CONFIG_FILE"
      StrCpy $HasExistingConfig "1"
      StrCpy $ConfigLoadPath "$CONFIG_FILE"
      StrCpy $InstallType "advanced"
      Call LoadExistingConfig
    ${EndIf}
    Call DetectNginx
  ${EndIf}

  ${If} ${Silent}
    StrCpy $OpenUI "0"
  ${EndIf}

  Call ParseArguments

  ${If} $ConfigImport != ""
    ${IfNot} ${FileExists} "$ConfigImport"
      ${If} ${FileExists} "$EXEDIR\$ConfigImport"
        StrCpy $ConfigImport "$EXEDIR\$ConfigImport"
      ${Else}
        MessageBox MB_ICONSTOP|MB_OK "$(NI_CONFIG_IMPORT_FAILED)"
        Abort
      ${EndIf}
    ${EndIf}
    StrCpy $ConfigLoadPath "$ConfigImport"
    Call LoadExistingConfig
    StrCpy $HasExistingConfig "1"
    StrCpy $InstallType "advanced"
    StrCpy $ConfigTouched "1"
  ${ElseIf} ${Silent}
    ${If} $DatabaseDriver == "postgres"
      ${If} $DbHost == ""
      ${OrIf} $DbName == ""
      ${OrIf} $DbUsername == ""
      ${OrIf} $DbPassword == ""
        MessageBox MB_ICONSTOP|MB_OK "$(NI_DB_REQUIRED)"
        Abort
      ${EndIf}
    ${EndIf}
    ${If} $ServiceAccountKind == "custom"
      ${If} $IsUpgrade == "1"
      ${AndIf} $ServiceAccountArgGiven != "1"
        ; Keeping the existing service account; the password is not needed.
      ${Else}
        ${If} $ServiceAccount == ""
        ${OrIf} $ServicePassword == ""
          MessageBox MB_ICONSTOP|MB_OK "$(NI_SERVICE_CUSTOM_REQUIRED)"
          Abort
        ${EndIf}
      ${EndIf}
    ${EndIf}
  ${EndIf}
FunctionEnd

Function WelcomePagePre
  StrCpy $WelcomeText "$(NI_WELCOME_TEXT)$\r$\n$\r$\n$(NI_VERSION_LABEL): ${VERSION}"
  ${If} $IsUpgrade == "1"
    StrCpy $WelcomeText "$WelcomeText$\r$\n$\r$\n$(NI_UPGRADE_TEXT)"
  ${EndIf}
FunctionEnd

;--------------------------------
;Helpers

Function SkipOnUpgrade
  ${If} $IsUpgrade == "1"
    Abort
  ${EndIf}
FunctionEnd

Function FinishPagePre
  Call BuildUiUrl
  StrCpy $FinishText "$(NI_FINISH_TEXT)$\r$\n$\r$\n$0"
FunctionEnd

Function BuildUiUrl
  ${If} $ServerAddress == "0.0.0.0"
  ${OrIf} $ServerAddress == ""
    StrCpy $0 "127.0.0.1"
  ${ElseIf} $ServerAddress == "::"
    StrCpy $0 "[::1]"
  ${Else}
    ${StrStr} $1 "$ServerAddress" ":"
    ${If} $1 != ""
      StrCpy $0 "[$ServerAddress]"
    ${Else}
      StrCpy $0 "$ServerAddress"
    ${EndIf}
  ${EndIf}
  StrCpy $0 "http://$0:$ServerPort/"
FunctionEnd

Function OpenUI
  Call BuildUiUrl
  ExecShell "open" "$0"
FunctionEnd

Function CreateShortcuts
  DetailPrint "$(NI_DETAIL_SHORTCUTS)"
  SetShellVarContext all

  Call BuildUiUrl
  StrCpy $UiUrl "$0"

  ${If} $StartMenuShortcut == "1"
    CreateDirectory "$SMPROGRAMS\${APP_NAME}"
    CreateShortCut "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk" \
      "$WINDIR\explorer.exe" "$\"$UiUrl$\"" \
      "$INSTDIR\nginx-ignition.ico" 0 SW_SHOWNORMAL "" "$(NI_SHORTCUT_NAME)"
    CreateShortCut "$SMPROGRAMS\${APP_NAME}\Uninstall ${APP_NAME}.lnk" \
      "$INSTDIR\uninstall.exe" "" "$INSTDIR\nginx-ignition.ico" 0 SW_SHOWNORMAL "" "$(NI_UNINSTALL_SHORTCUT)"
  ${Else}
    Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
    Delete "$SMPROGRAMS\${APP_NAME}\Uninstall ${APP_NAME}.lnk"
    RMDir "$SMPROGRAMS\${APP_NAME}"
  ${EndIf}

  ${If} $DesktopShortcut == "1"
    CreateShortCut "$DESKTOP\${APP_NAME}.lnk" \
      "$WINDIR\explorer.exe" "$\"$UiUrl$\"" \
      "$INSTDIR\nginx-ignition.ico" 0 SW_SHOWNORMAL "" "$(NI_SHORTCUT_NAME)"
  ${Else}
    Delete "$DESKTOP\${APP_NAME}.lnk"
  ${EndIf}
FunctionEnd

Function WriteRegistry
  WriteRegStr HKLM "${APP_REGKEY}" "InstallDir" "$INSTDIR"
  WriteRegStr HKLM "${APP_REGKEY}" "ConfigFile" "$CONFIG_FILE"
  WriteRegStr HKLM "${APP_REGKEY}" "DataDir" "$DataDir"
  WriteRegStr HKLM "${APP_REGKEY}" "DatabaseDriver" "$DatabaseDriver"
  WriteRegStr HKLM "${APP_REGKEY}" "ServiceAccountKind" "$ServiceAccountKind"
  WriteRegStr HKLM "${APP_REGKEY}" "ServiceAccount" "$ServiceAccount"
  WriteRegStr HKLM "${APP_REGKEY}" "ServiceStartType" "$ServiceStartType"
  WriteRegStr HKLM "${APP_REGKEY}" "FirewallUI" "$FirewallUI"
  WriteRegStr HKLM "${APP_REGKEY}" "FirewallHTTP" "$FirewallHTTP"
  WriteRegStr HKLM "${APP_REGKEY}" "FirewallHTTPS" "$FirewallHTTPS"
  WriteRegStr HKLM "${APP_REGKEY}" "DesktopShortcut" "$DesktopShortcut"
  WriteRegStr HKLM "${APP_REGKEY}" "StartMenuShortcut" "$StartMenuShortcut"
  WriteRegStr HKLM "${APP_REGKEY}" "Version" "${VERSION}"

  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "DisplayName" "${APP_NAME}"
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "DisplayVersion" "${VERSION}"
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "DisplayIcon" "$INSTDIR\nginx-ignition.ico"
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "Publisher" "Lucas Dillmann"
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "URLInfoAbout" "${APP_URL}"
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "HelpLink" "${APP_URL}"
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegStr HKLM "${APP_UNINSTALLKEY}" "QuietUninstallString" '"$INSTDIR\uninstall.exe" /S'
  WriteRegDWORD HKLM "${APP_UNINSTALLKEY}" "NoModify" 1
  WriteRegDWORD HKLM "${APP_UNINSTALLKEY}" "NoRepair" 1
FunctionEnd

;--------------------------------
;Installer

Section "$(NI_SECTION_APP)"

  SetShellVarContext all

  Call StopService

  SetOutPath "$INSTDIR"
  IfErrors dir_error

  File "/oname=nginx-ignition.exe" "${BINARY}"
  File "/oname=nginx-ignition.ico" "${ASSETSDIR}/nginx-ignition.ico"
  File "/oname=LICENSE.txt" "${ROOT}/LICENSE"
  IfErrors dir_error
  WriteUninstaller "$INSTDIR\uninstall.exe"

  ${If} $ConfigImport != ""
    StrCpy $ConfigLoadPath "$ConfigImport"
    Call LoadExistingConfig
  ${EndIf}

  Call EnsureDirectories
  Call WriteConfiguration
  Call ApplyPermissions
  Call RegisterEventLog
  Call ConfigureService
  Call ApplyFirewall
  Call CreateShortcuts
  Call WriteRegistry
  Call StartService

  !insertmacro MUI_LANGDLL_SAVELANGUAGE

  ${If} ${Silent}
  ${AndIf} $OpenUI == "1"
    Call OpenUI
  ${EndIf}

  Goto section_done

  dir_error:
    MessageBox MB_ICONSTOP|MB_OK "$(NI_ERR_DIR)"
    Abort

  section_done:

SectionEnd

;--------------------------------
;Uninstaller

Function un.onInit
  !insertmacro MUI_UNGETLANGUAGE

  SetShellVarContext all

  StrCpy $APPDATA_DIR "$APPDATA\${APP_ID}"
  StrCpy $UninstallPurge "0"

  ReadRegStr $0 HKLM "${APP_REGKEY}" "InstallDir"
  ${If} $0 != ""
    StrCpy $INSTDIR $0
  ${EndIf}

  ReadRegStr $0 HKLM "${APP_REGKEY}" "ConfigFile"
  ${If} $0 != ""
    StrCpy $CONFIG_FILE $0
  ${EndIf}

  ReadRegStr $0 HKLM "${APP_REGKEY}" "DataDir"
  ${If} $0 != ""
    StrCpy $DataDir $0
  ${EndIf}

  InitPluginsDir
  File "/oname=$PLUGINSDIR\wait-service.ps1" "${SCRIPTDIR}/wait-service.ps1"
FunctionEnd

;--------------------------------
;Uninstaller options

Function un.OptionsPageCreate
  !insertmacro MUI_HEADER_TEXT "$(NI_UNINSTALL_OPTIONS_TITLE)" "$(NI_UNINSTALL_OPTIONS_SUBTITLE)"

  nsDialogs::Create 1018
  Pop $0
  ${If} $0 == error
    Abort
  ${EndIf}

  ${NSD_CreateGroupBox} 0 0 300u 62u "$(NI_GROUP_DATA)"
  Pop $0

  ${NSD_CreateCheckbox} 10u 14u 280u 10u "$(NI_UNINSTALL_PURGE_CHECKBOX)"
  Pop $hUninstallPurge
  ${If} $UninstallPurge == "1"
    ${NSD_Check} $hUninstallPurge
  ${EndIf}

  ${NSD_CreateLabel} 10u 30u 280u 26u "$(NI_UNINSTALL_PURGE_HINT)"
  Pop $0

  nsDialogs::Show
FunctionEnd

Function un.OptionsPageLeave
  ${NSD_GetState} $hUninstallPurge $UninstallPurge
FunctionEnd

Function un.StopService
  ClearErrors
  ReadRegStr $0 HKLM "SYSTEM\CurrentControlSet\Services\${APP_SERVICE}" "ImagePath"
  ${If} ${Errors}
  ${OrIf} $0 == ""
    Return
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
    nsExec::ExecToStack 'taskkill.exe /F /IM nginx-ignition.exe'
    Pop $0
    Pop $1
  ${EndIf}

  nsExec::ExecToStack 'sc.exe delete ${APP_SERVICE}'
  Pop $0
  Pop $1
FunctionEnd

Function un.RemoveFirewall
  nsExec::ExecToStack 'netsh advfirewall firewall delete rule name="${APP_RULE_WEB}"'
  nsExec::ExecToStack 'netsh advfirewall firewall delete rule name="${APP_RULE_HTTP}"'
  nsExec::ExecToStack 'netsh advfirewall firewall delete rule name="${APP_RULE_HTTPS}"'
  Pop $0
  Pop $0
  Pop $0
  Pop $0
  Pop $0
  Pop $0
FunctionEnd

Section "Uninstall"

  SetShellVarContext all

  DetailPrint "$(NI_DETAIL_UNINSTALL)"

  Call un.StopService
  Call un.RemoveFirewall

  DeleteRegKey HKLM "SYSTEM\CurrentControlSet\Services\EventLog\Application\${APP_SERVICE}"
  DeleteRegKey HKLM "${APP_UNINSTALLKEY}"
  DeleteRegKey HKLM "${APP_REGKEY}"

  Delete "$SMPROGRAMS\${APP_NAME}\${APP_NAME}.lnk"
  Delete "$SMPROGRAMS\${APP_NAME}\Uninstall ${APP_NAME}.lnk"
  RMDir "$SMPROGRAMS\${APP_NAME}"
  Delete "$DESKTOP\${APP_NAME}.lnk"

  Delete "$INSTDIR\nginx-ignition.exe"
  Delete "$INSTDIR\nginx-ignition.ico"
  Delete "$INSTDIR\LICENSE.txt"
  Delete "$INSTDIR\uninstall.exe"
  RMDir "$INSTDIR"

  ${If} ${Silent}
    StrCpy $R1 "0"
    ${GetParameters} $R0
    ${GetOptions} $R0 "/PURGE=" $R1
    StrCmp $R1 "1" 0 un_keep_data
      RMDir /r "$APPDATA_DIR"
      Goto un_data_done
  ${Else}
    StrCmp $UninstallPurge "1" 0 un_keep_data
      RMDir /r "$APPDATA_DIR"
      Goto un_data_done
  ${EndIf}

  un_keep_data:
    DetailPrint "$(NI_UNINSTALL_KEEP_DATA)"

  un_data_done:

SectionEnd
