@echo off

set BUILD_32=false

set SRC_DIR=%USERPROFILE%\Desktop\projects\Hypothermia\src
set BUILD_DIR=%USERPROFILE%\Desktop\projects\Hypothermia\build

cd /d "%SRC_DIR%" 2>nul
if %errorlevel% neq 0 (
	cd /d "%SRC_DIR%"
)

if not exist "%BUILD_DIR%" (
	mkdir "%BUILD_DIR%"
)

set GOARCH=amd64

cd /d "%SRC_DIR%"

go build -trimpath -ldflags="-w -s -H=windowsgui" -o "%BUILD_DIR%\Hypothermia.exe" main.go
if %errorlevel% neq 0 (
	color 0C
	echo Hypothermia build failed

	exit /b %errorlevel%
)

if /i "%BUILD_32%"=="true" (
	setlocal
	set GOARCH=386

	go build -trimpath -ldflags="-w -s -H=windowsgui" -o "%BUILD_DIR%\Hypothermia_x86.exe" main.go
	if %errorlevel% neq 0 (
		color 0C
		echo Hypothermia 32-bit build failed
		endlocal
		exit /b %errorlevel%
	)

	endlocal
)

color 0A
echo Hypothermia built successfully
cls

color 07
