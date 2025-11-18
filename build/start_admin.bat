@echo off
REM 启动所有API服务
setlocal enabledelayedexpansion

REM 获取脚本所在目录的上级目录（beaver-server根目录）
for %%i in ("%~dp0..") do set "ROOT_DIR=%%~fi"

REM 启动所有Admin服务
wt new-tab --title "Gateway Admin" cmd /k "cd /d !ROOT_DIR!\app\gateway\gateway_admin && go run gateway.go" ^
  ; new-tab --title "Backend Admin" cmd /k "cd /d !ROOT_DIR!\app\backend\backend_admin && go run backend.go" ^

