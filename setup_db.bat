@echo off
echo ==============================================
echo Wander: Local Guide Hub - Database Setup
echo ==============================================
echo.
set /p DB_SUPERUSER="Enter PostgreSQL superuser username (default: postgres): "
if "%DB_SUPERUSER%"=="" set DB_SUPERUSER=postgres

set /p APP_DB_PASS="Enter password for new database user 'wander_user' (default: Josedg0212): "
if "%APP_DB_PASS%"=="" set APP_DB_PASS=Josedg0212

echo.
echo Syncing password to configuration files...
powershell -Command "(gc backend\migrations\db_setup.sql) -replace 'PASSWORD ''[^'']+''', 'PASSWORD ''%APP_DB_PASS%''' | Out-File -encoding ascii backend\migrations\db_setup.sql"
powershell -Command "(gc backend\.env) -replace 'DB_PASSWORD=.*', 'DB_PASSWORD=%APP_DB_PASS%' | Out-File -encoding ascii backend\.env"

echo.
echo Please enter the password for the superuser '%DB_SUPERUSER%' when prompted.
psql -U %DB_SUPERUSER% -h localhost -d postgres -f backend\migrations\db_setup.sql
if %ERRORLEVEL% equ 0 (
    echo.
    echo Database 'wander_db' and user 'wander_user' configured successfully!
) else (
    echo.
    echo Failed to configure database. Please verify your credentials and psql path.
)
pause
