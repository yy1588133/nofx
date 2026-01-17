@echo off
setlocal EnableExtensions EnableDelayedExpansion

rem NOFX - Windows 一键重建并部署到 Docker
rem
rem 默认行为（双击运行）：
rem   - 使用 docker-compose.yml（本地源码构建）
rem   - 重新构建镜像并启动/更新容器（后台运行）
rem
rem 可选用法（命令行）：
rem   docker-redeploy.bat dev
rem   docker-redeploy.bat dev --no-cache
rem   docker-redeploy.bat stable
rem   docker-redeploy.bat prod
rem   docker-redeploy.bat dev --dry-run

chcp 65001 >nul

set "MODE=dev"
set "NO_CACHE="
set "DO_PULL="
set "DO_DOWN="
set "DRY_RUN="

:parseArgs
if "%~1"=="" goto :argsDone

if /I "%~1"=="dev" (
  set "MODE=dev"
  shift
  goto :parseArgs
)
if /I "%~1"=="stable" (
  set "MODE=stable"
  shift
  goto :parseArgs
)
if /I "%~1"=="prod" (
  set "MODE=prod"
  shift
  goto :parseArgs
)

if /I "%~1"=="--no-cache" (
  set "NO_CACHE=1"
  shift
  goto :parseArgs
)
if /I "%~1"=="--pull" (
  set "DO_PULL=1"
  shift
  goto :parseArgs
)
if /I "%~1"=="--down" (
  set "DO_DOWN=1"
  shift
  goto :parseArgs
)
if /I "%~1"=="--dry-run" (
  set "DRY_RUN=1"
  shift
  goto :parseArgs
)

if /I "%~1"=="-h" goto :help
if /I "%~1"=="--help" goto :help

echo [错误] 未识别参数: %~1
echo.
goto :help

:argsDone

rem 确保在脚本所在目录（项目根目录）执行
pushd "%~dp0" >nul 2>&1
if !errorlevel! neq 0 (
  echo [错误] 无法切换到脚本目录: "%~dp0"
  goto :fail_no_popd
)

set "COMPOSE_FILE=docker-compose.yml"
if /I "%MODE%"=="stable" set "COMPOSE_FILE=docker-compose.stable.yml"
if /I "%MODE%"=="prod" set "COMPOSE_FILE=docker-compose.prod.yml"

if not exist "%COMPOSE_FILE%" (
  echo [错误] 未找到 Compose 文件: "%COMPOSE_FILE%"
  echo [提示] 请在项目根目录运行，或检查文件是否存在。
  goto :fail
)

if not exist ".env" (
  echo [警告] 未找到 .env 文件，容器启动可能缺少必要环境变量。
)

rem 检测 compose 命令（优先 docker compose）
set "COMPOSE_CMD="
docker compose version >nul 2>&1
if !errorlevel! equ 0 set "COMPOSE_CMD=docker compose"
if not defined COMPOSE_CMD (
  docker-compose version >nul 2>&1
  if !errorlevel! equ 0 set "COMPOSE_CMD=docker-compose"
)

if not defined COMPOSE_CMD (
  if defined DRY_RUN (
    set "COMPOSE_CMD=docker compose"
  ) else (
    echo [错误] 未检测到 Docker Compose（docker compose / docker-compose）。
    echo [提示] 请确认已安装 Docker Desktop，并启用 Docker Compose。
    goto :fail
  )
)

rem stable/prod 默认拉取最新镜像
if /I "%MODE%"=="stable" if not defined DO_PULL set "DO_PULL=1"
if /I "%MODE%"=="prod" if not defined DO_PULL set "DO_PULL=1"

echo [信息] 模式: %MODE%
echo [信息] Compose: %COMPOSE_FILE%
echo [信息] 命令: %COMPOSE_CMD%
echo.

if defined DRY_RUN (
  echo [DRY-RUN] %COMPOSE_CMD% -f "%COMPOSE_FILE%" ^(以下为将执行的命令^)
  if defined DO_DOWN echo [DRY-RUN] %COMPOSE_CMD% -f "%COMPOSE_FILE%" down --remove-orphans
  if /I "%MODE%"=="dev" (
    set "BUILD_ARGS="
    if defined NO_CACHE set "BUILD_ARGS=!BUILD_ARGS! --no-cache"
    if defined DO_PULL set "BUILD_ARGS=!BUILD_ARGS! --pull"
    echo [DRY-RUN] %COMPOSE_CMD% -f "%COMPOSE_FILE%" build!BUILD_ARGS!
  ) else (
    if defined DO_PULL echo [DRY-RUN] %COMPOSE_CMD% -f "%COMPOSE_FILE%" pull
  )
  echo [DRY-RUN] %COMPOSE_CMD% -f "%COMPOSE_FILE%" up -d --remove-orphans
  echo [DRY-RUN] %COMPOSE_CMD% -f "%COMPOSE_FILE%" ps
  goto :done
)

rem 检查 docker 是否可用（仅非 dry-run）
where docker >nul 2>&1
if !errorlevel! neq 0 (
  echo [错误] 未找到 docker 命令，请先安装 Docker Desktop。
  goto :fail
)

docker info >nul 2>&1
if !errorlevel! neq 0 (
  echo [错误] Docker 守护进程未就绪。
  echo [提示] 请确认 Docker Desktop 已启动，并处于 Linux Containers 模式。
  goto :fail
)

if defined DO_DOWN (
  echo [步骤] 停止并移除现有容器（不删除 volumes）...
  %COMPOSE_CMD% -f "%COMPOSE_FILE%" down --remove-orphans
  if !errorlevel! neq 0 goto :fail
  echo.
)

if /I "%MODE%"=="dev" (
  echo [步骤] 构建镜像...
  set "BUILD_ARGS="
  if defined NO_CACHE set "BUILD_ARGS=!BUILD_ARGS! --no-cache"
  if defined DO_PULL set "BUILD_ARGS=!BUILD_ARGS! --pull"
  %COMPOSE_CMD% -f "%COMPOSE_FILE%" build!BUILD_ARGS!
  if !errorlevel! neq 0 goto :fail
  echo.
) else (
  if defined DO_PULL (
    echo [步骤] 拉取镜像...
    %COMPOSE_CMD% -f "%COMPOSE_FILE%" pull
    if !errorlevel! neq 0 goto :fail
    echo.
  )
)

echo [步骤] 启动/更新容器...
%COMPOSE_CMD% -f "%COMPOSE_FILE%" up -d --remove-orphans
if !errorlevel! neq 0 goto :fail
echo.

echo [步骤] 当前容器状态:
%COMPOSE_CMD% -f "%COMPOSE_FILE%" ps
if !errorlevel! neq 0 goto :fail

echo.
echo [完成] 部署完成。可访问:
echo   - 前端: http://127.0.0.1:3000
echo   - 后端: http://127.0.0.1:8080

:done
popd >nul 2>&1
exit /b 0

:help
echo 用法:
echo   docker-redeploy.bat [dev^|stable^|prod] [--no-cache] [--pull] [--down] [--dry-run]
echo.
echo 说明:
echo   dev    : 使用 docker-compose.yml（本地源码构建，默认）
echo   stable : 使用 docker-compose.stable.yml（稳定镜像）
echo   prod   : 使用 docker-compose.prod.yml（latest 镜像）
echo.
echo 选项:
echo   --no-cache  dev 模式构建不使用缓存
echo   --pull      dev 模式构建时拉取基础镜像；stable/prod 拉取服务镜像
echo   --down      先执行 down（不删除 volumes）
echo   --dry-run   仅打印将要执行的命令，不实际执行
echo.
exit /b 0

:fail
echo.
echo [失败] 部署未完成，请根据上方输出排查。
popd >nul 2>&1
pause
exit /b 1

:fail_no_popd
echo.
echo [失败] 部署未完成，请根据上方输出排查。
pause
exit /b 1

