@echo off
REM 启动所有 API 服务（track/dictionary 已并入 platform_api）
setlocal enabledelayedexpansion

for %%i in ("%~dp0..") do set "ROOT_DIR=%%~fi"

wt new-tab --title "Gateway API" cmd /k "cd /d !ROOT_DIR!\app\gateway\gateway_api && go run gateway.go" ^
  ; new-tab --title "Auth API" cmd /k "cd /d !ROOT_DIR!\app\auth\auth_api && go run auth.go" ^
  ; new-tab --title "User API" cmd /k "cd /d !ROOT_DIR!\app\user\user_api && go run user.go" ^
  ; new-tab --title "Friend API" cmd /k "cd /d !ROOT_DIR!\app\friend\friend_api && go run friend.go" ^
  ; new-tab --title "Group API" cmd /k "cd /d !ROOT_DIR!\app\group\group_api && go run group.go" ^
  ; new-tab --title "Chat API" cmd /k "cd /d !ROOT_DIR!\app\chat\chat_api && go run chat.go" ^
  ; new-tab --title "WS API" cmd /k "cd /d !ROOT_DIR!\app\ws\ws_api && go run ws.go" ^
  ; new-tab --title "File API" cmd /k "cd /d !ROOT_DIR!\app\file\file_api && go run file.go" ^
  ; new-tab --title "Emoji API" cmd /k "cd /d !ROOT_DIR!\app\emoji\emoji_api && go run emoji.go" ^
  ; new-tab --title "Moment API" cmd /k "cd /d !ROOT_DIR!\app\moment\moment_api && go run moment.go" ^
  ; new-tab --title "Platform API" cmd /k "cd /d !ROOT_DIR!\app\platform\platform_api && go run platform.go" ^
  ; new-tab --title "Datasync API" cmd /k "cd /d !ROOT_DIR!\app\datasync\datasync_api && go run datasync.go" ^
  ; new-tab --title "Notification API" cmd /k "cd /d !ROOT_DIR!\app\notification\notification_api && go run notification.go" ^
  ; new-tab --title "Call API" cmd /k "cd /d !ROOT_DIR!\app\call\call_api && go run call.go" ^
  ; new-tab --title "Open API" cmd /k "cd /d !ROOT_DIR!\app\open\open_api && go run open.go"
