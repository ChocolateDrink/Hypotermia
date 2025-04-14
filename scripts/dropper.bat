@echo off

set HELPERS_DIR=%USERPROFILE%\Desktop\projects\Hypothermia\_helpers
set BUILD_DIR=%USERPROFILE%\Desktop\projects\Hypothermia\build

if not exist "%BUILD_DIR%" (
    mkdir "%BUILD_DIR%"
)

cd /d "%HELPERS_DIR%"

go build -trimpath -ldflags="-w -s" -o "%BUILD_DIR%\Dropper.exe" dropper.go
if %errorlevel% neq 0 (
    color 0C
    echo Dropper build failed

    exit /b %errorlevel%
)

color 0A
echo Dropper built successfully
cls

color 07