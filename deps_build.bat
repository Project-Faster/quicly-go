
@echo off

echo *********************************************
echo *****    QUICLY-GO DEPENDENCIES BUILD   *****
echo *********************************************
echo.

set BASEDIR=%~dp0
echo %BASEDIR%

go env GOOS > tmpfile
set /p GOOS= < tmpfile

go env GOARCH > tmpfile
set /p GOARCH= < tmpfile
del tmpfile

echo [%GOOS%]
echo [%GOARCH%]

set /A BUILD="Release"

if "%1" EQU "--quic" (
    goto quic
)

if "%1" EQU "--clean" (
    goto reset
)
if "%2" EQU "--clean" (
    goto reset
)

if "%1" EQU "--debug" (
    set /A BUILD="Debug"
)

if "%1" EQU "--help" (
    echo "Usage: build.bat [--debug][--clean]"
    exit /B 1
)

ECHO [Prerequisites check: GO]
go version
if %ERRORLEVEL% GEQ 1 goto fail

ECHO [Prerequisites check: cmake]
cmake --version
if %ERRORLEVEL% GEQ 1 goto fail

ECHO [Prerequisites check: git]
git version
if %ERRORLEVEL% GEQ 1 goto fail

echo [Init submodules]
git submodule init deps/openssl
git submodule init deps/quicly
git submodule init deps/c-for-go

echo [Update submodules]
git submodule update --init --recursive

echo [Create gen directories]
mkdir gen_openssl
mkdir gen_quicly

echo [Build OpenSSL]
pushd gen_openssl
cmake ../deps/openssl -G"MinGW Makefiles" -DCMAKE_INSTALL_PREFIX=%BASEDIR%/internal/%GOOS%/%GOARCH%/bindings -DCMAKE_BUILD_TYPE=%BUILD%
if %ERRORLEVEL% NEQ 0 goto fail

cmake --build .
if %ERRORLEVEL% NEQ 0 goto fail

cmake --build . --target install
if %ERRORLEVEL% NEQ 0 goto fail
popd

:quic
echo [Build Quicly]
pushd gen_quicly
cmake ../deps/quicly -G"MinGW Makefiles" -DCMAKE_INSTALL_PREFIX=%BASEDIR%/internal/%GOOS%/%GOARCH%/bindings -DOPENSSL_ROOT_DIR=%BASEDIR%/%GOOS%/%GOARCH%/bindings/include ^
                                         -DCMAKE_BUILD_TYPE=%BUILD% -DWITH_EXAMPLE=OFF -DCMAKE_VERBOSE_MAKEFILE=ON
if %ERRORLEVEL% NEQ 0 goto fail

cmake --build .
if %ERRORLEVEL% NEQ 0 goto fail

cmake --build . --target install
if %ERRORLEVEL% NEQ 0 goto fail
popd

:ok
echo.
echo ********************************
echo **** RESULT: SUCCESS        ****
echo ********************************
cd %BASEDIR%
exit /B 0

:fail
echo.
echo ********************************
echo **** RESULT: FAILURE        ****
echo ********************************
cd %BASEDIR%
exit /B 1

:reset
cd %BASEDIR%
go clean -cache -x
git reset --hard
git clean -f -d
exit /B 0
