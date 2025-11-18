@echo off
REM 启动所有API服务
setlocal enabledelayedexpansion

REM 获取脚本所在目录的上级目录（beaver-server根目录）
for %%i in ("%~dp0..") do set "ROOT_DIR=%%~fi"

wt new-tab --title "Auth API" cmd /k "cd /d !ROOT_DIR!\app\auth\auth_api && go run auth.go" ^
  ; new-tab --title "User API" cmd /k "cd /d !ROOT_DIR!\app\user\user_api && go run user.go" ^
  ; new-tab --title "File API" cmd /k "cd /d !ROOT_DIR!\app\file\file_api && go run file.go" ^
  ; new-tab --title "Friend API" cmd /k "cd /d !ROOT_DIR!\app\friend\friend_api && go run friend.go" ^
  ; new-tab --title "Chat API" cmd /k "cd /d !ROOT_DIR!\app\chat\chat_api && go run chat.go" ^
  ; new-tab --title "WS API" cmd /k "cd /d !ROOT_DIR!\app\ws\ws_api && go run ws.go" ^
  ; new-tab --title "Group API" cmd /k "cd /d !ROOT_DIR!\app\group\group_api && go run group.go" ^
  ; new-tab --title "Gateway API" cmd /k "cd /d !ROOT_DIR!\app\gateway\gateway_api && go run gateway.go" ^
  ; new-tab --title "Emoji API" cmd /k "cd /d !ROOT_DIR!\app\emoji\emoji_api && go run emoji.go" ^
  ; new-tab --title "Moment API" cmd /k "cd /d !ROOT_DIR!\app\moment\moment_api && go run moment.go" ^
  ; new-tab --title "Dictionary API" cmd /k "cd /d !ROOT_DIR!\app\dictionary\dictionary_api && go run dictionary.go" ^
  ; new-tab --title "Track API" cmd /k "cd /d !ROOT_DIR!\app\track\track_api && go run track.go" ^
  ; new-tab --title "Update API" cmd /k "cd /d !ROOT_DIR!\app\update\update_api && go run update.go" ^
  ; new-tab --title "Datasync API" cmd /k "cd /d !ROOT_DIR!\app\datasync\datasync_api && go run datasync.go" 