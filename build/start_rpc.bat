@echo off
REM 启动所有RPC服务
setlocal enabledelayedexpansion

REM 获取脚本所在目录的上级目录（beaver-server根目录）
for %%i in ("%~dp0..") do set "ROOT_DIR=%%~fi"

wt new-tab --title "User RPC" cmd /k "cd /d !ROOT_DIR!\app\user\user_rpc && go run userrpc.go" ^
  ; new-tab --title "Group RPC" cmd /k "cd /d !ROOT_DIR!\app\group\group_rpc && go run grouprpc.go" ^
  ; new-tab --title "Friend RPC" cmd /k "cd /d !ROOT_DIR!\app\friend\friend_rpc && go run friendrpc.go" ^
  ; new-tab --title "Chat RPC" cmd /k "cd /d !ROOT_DIR!\app\chat\chat_rpc && go run chatrpc.go" ^
  ; new-tab --title "File RPC" cmd /k "cd /d !ROOT_DIR!\app\file\file_rpc && go run filerpc.go" ^
  ; new-tab --title "Dictionary RPC" cmd /k "cd /d !ROOT_DIR!\app\dictionary\dictionary_rpc && go run dictionaryrpc.go" ^
  ; new-tab --title "Emoji RPC" cmd /k "cd /d !ROOT_DIR!\app\emoji\emoji_rpc && go run emojirpc.go" ^
  ; new-tab --title "Notification RPC" cmd /k "cd /d !ROOT_DIR!\app\notification\notification_rpc && go run notificationrpc.go" ^