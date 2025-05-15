@echo off

REM Set environment variables
setlocal enabledelayedexpansion

REM Check if .env file exists
if not exist ".env" (
    echo Error: .env file does not exist!
    echo Please create .env file and set necessary environment variables.
    echo Example:
    echo # QWeather API Configuration
    echo QWEATHER_API_BASE=https://api.qweather.com
    echo QWEATHER_API_KEY=your_api_key_here
    pause
    exit /b 1
)

REM Read environment variables from .env file
for /F "usebackq tokens=1,* delims==" %%G in (".env") do (
    set "line=%%G"
    if not "!line:~0,1!"=="#" (
        set "%%G=%%H"
    )
)

REM Display set environment variables
echo Environment variables set:
echo QWEATHER_API_BASE=%QWEATHER_API_BASE%
echo QWEATHER_API_KEY=%QWEATHER_API_KEY%

REM Check if necessary environment variables are set
if "%QWEATHER_API_BASE%"=="" (
    echo Error: QWEATHER_API_BASE environment variable not set!
    pause
    exit /b 1
)

if "%QWEATHER_API_KEY%"=="" (
    echo Error: QWEATHER_API_KEY environment variable not set!
    pause
    exit /b 1
)

REM Run the program
echo Starting QWeather MCP server...
go run main.go

endlocal
