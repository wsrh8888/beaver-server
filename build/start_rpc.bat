@echo off
REM 启动所有RPC服务
wt new-tab --title "User RPC" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\user\user_rpc && go run userrpc.go" ^
  ; new-tab --title "Group RPC" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\group\group_rpc && go run grouprpc.go" ^
  ; new-tab --title "Friend RPC" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\friend\friend_rpc && go run friendrpc.go" ^
  ; new-tab --title "Chat RPC" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\chat\chat_rpc && go run chatrpc.go" ^
  ; new-tab --title "File RPC" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\file\file_rpc && go run filerpc.go" ^
  ; new-tab --title "Dictionary RPC" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\dictionary\dictionary_rpc && go run dictionaryrpc.go" ^
  ; new-tab --title "Datasync RPC" cmd /k "cd /d F:\code\mine\IM\beaver-server\app\datasync\datasync_rpc && go run datasyncrpc.go"