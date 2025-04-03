@echo off

set SRC_DIR=%USERPROFILE%\Desktop\projects\Hypotermia\src
set BUILD_DIR=%USERPROFILE%\Desktop\projects\Hypotermia\build

cd /d "%SRC_DIR%" 2>nul
if %errorlevel% neq 0 (
    cd /d "%SRC_DIR%"
)

if not exist "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%"
)

cd /d "%SRC_DIR%"
go build -ldflags -H=windowsgui -o ../build/Hypotermia.exe main.go
if %errorlevel% neq 0 (
    color 0C
    echo Hypotermia build failed

    exit /b %errorlevel%
)

go build -o ../build/HypotermiaDebug.exe main.go
if %errorlevel% neq 0 (
    color 0C
    echo Hypotermia debug build failed

    exit /b %errorlevel%
)

color 0A
echo Hypotermia built successfully
cls

color 07
