@echo off

set SRC_DIR=%USERPROFILE%\Desktop\projects\Hypothermia\src
set BUILD_DIR=%USERPROFILE%\Desktop\projects\Hypothermia\build

cd /d "%SRC_DIR%" 2>nul
if %errorlevel% neq 0 (
    cd /d "%SRC_DIR%"
)

if not exist "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%"
)

cd /d "%SRC_DIR%"
go build -ldflags -H=windowsgui -o ../build/Hypothermia.exe main.go
if %errorlevel% neq 0 (
    color 0C
    echo Hypothermia build failed

    exit /b %errorlevel%
)

go build -o ../build/HypothermiaDebug.exe main.go
if %errorlevel% neq 0 (
    color 0C
    echo Hypothermia debug build failed

    exit /b %errorlevel%
)

color 0A
echo Hypothermia built successfully
cls

color 07