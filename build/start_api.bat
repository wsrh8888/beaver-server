@echo off
REM 启动所有 API 服务（user+auth->account，file+emoji->asset，chat->message，friend+group->relation）
setlocal enabledelayedexpansion

for %%i in ("%~dp0..") do set "ROOT_DIR=%%~fi"

wt new-tab --title "Gateway API" cmd /k "cd /d !ROOT_DIR!\app\gateway\gateway_api && go run gateway.go" ^
  ; new-tab --title "Account API" cmd /k "cd /d !ROOT_DIR!\app\account\account_api && go run account.go" ^
  ; new-tab --title "Asset API" cmd /k "cd /d !ROOT_DIR!\app\asset\asset_api && go run asset.go" ^
  ; new-tab --title "Relation API" cmd /k "cd /d !ROOT_DIR!\app\relation\relation_api && go run relation.go" ^
  ; new-tab --title "Message API" cmd /k "cd /d !ROOT_DIR!\app\message\message_api && go run message.go" ^
  ; new-tab --title "Platform API" cmd /k "cd /d !ROOT_DIR!\app\platform\platform_api && go run platform.go" ^
  ; new-tab --title "Notification API" cmd /k "cd /d !ROOT_DIR!\app\notification\notification_api && go run notification.go" ^
  ; new-tab --title "Datasync API" cmd /k "cd /d !ROOT_DIR!\app\datasync\datasync_api && go run datasync.go" ^
  ; new-tab --title "WS API" cmd /k "cd /d !ROOT_DIR!\app\ws\ws_api && go run ws.go"
