@echo off
REM 启动所有Admin服务
wt new-tab --title "Auth Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\auth\auth_admin && go run auth.go" ^
  ; new-tab --title "Gateway Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\gateway\gateway_admin && go run gateway.go" ^
  ; new-tab --title "System Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\system\system_admin && go run system.go" ^
  ; new-tab --title "User Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\user\user_admin && go run user.go" ^
  ; new-tab --title "Group Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\group\group_admin && go run group.go" ^
  ; new-tab --title "Chat Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\chat\chat_admin && go run chat.go" ^
  ; new-tab --title "File Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\file\file_admin && go run file.go" ^
  ; new-tab --title "Friend Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\friend\friend_admin && go run friend.go" ^
  ; new-tab --title "Moment Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\moment\moment_admin && go run moment.go" ^
  ; new-tab --title "Feedback Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\feedback\feedback_admin && go run feedback.go" ^
  ; new-tab --title "Emoji Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\emoji\emoji_admin && go run emoji.go" ^
  ; new-tab --title "Update Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\update\update_admin && go run update.go" ^
  ; new-tab --title "Track Admin" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\track\track_admin && go run track.go" 