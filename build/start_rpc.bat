@echo off
REM 启动所有 RPC 服务（user+auth->account，file+emoji->asset，chat->message，friend+group->relation）
setlocal enabledelayedexpansion

for %%i in ("%~dp0..") do set "ROOT_DIR=%%~fi"

wt new-tab --title "Account RPC" cmd /k "cd /d !ROOT_DIR!\app\account\account_rpc && go run accountrpc.go" ^
  ; new-tab --title "Asset RPC" cmd /k "cd /d !ROOT_DIR!\app\asset\asset_rpc && go run assetrpc.go" ^
  ; new-tab --title "Relation RPC" cmd /k "cd /d !ROOT_DIR!\app\relation\relation_rpc && go run relationrpc.go" ^
  ; new-tab --title "Message RPC" cmd /k "cd /d !ROOT_DIR!\app\message\message_rpc && go run messagerpc.go" ^
  ; new-tab --title "Platform RPC" cmd /k "cd /d !ROOT_DIR!\app\platform\platform_rpc && go run platformrpc.go" ^
  ; new-tab --title "Notification RPC" cmd /k "cd /d !ROOT_DIR!\app\notification\notification_rpc && go run notificationrpc.go"
