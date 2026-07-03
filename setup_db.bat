@echo off
echo ==============================================
echo Wander: Local Guide Hub - Database Setup
echo ==============================================
echo.
set /p DB_SUPERUSER="Enter PostgreSQL superuser username (default: postgres): "
if "%DB_SUPERUSER%"=="" set DB_SUPERUSER=postgres
echo.
echo Please enter the password for the superuser '%DB_SUPERUSER%' when prompted.
psql -U %DB_SUPERUSER% -h localhost -f backend\migrations\db_setup.sql
if %ERRORLEVEL% equ 0 (
    echo.
    echo Database 'wander_db' and user 'wander_user' configured successfully!
) else (
    echo.
    echo Failed to configure database. Please verify your credentials and psql path.
)
pause
