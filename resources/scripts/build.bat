@echo off
setlocal enabledelayedexpansion
rem Generate and build the Showtunes API server code and executable

rem Set executable name
set execname=showtunesapi.exe

rem Set absolute project path
for %%a in ("%~dp0\.") do for %%b in ("%%~dpa\.") do for %%c in ("%%~dpb\.") do (
    set projectpath=%%~dpnxc
)
echo Project path: %projectpath%

rem Set executable build directory
set buildpath=%projectpath%\build

rem Make the build directory
if not exist %buildpath% (
    md %buildpath%
)
echo Build path: %buildpath%

rem Build the project file
go build -o %execname% .\cmd\main.go
if %errorlevel% neq 0 (
    echo Failed to compile main package
    exit /b %errorlevel%
)
rem Remove existing binary
del %buildpath%\*.exe 1>NUL 2>NUL
rem Move the new binary
move %execname% %buildpath%\ 1>NUL
echo Project built successfully to %buildpath%\%execname%
