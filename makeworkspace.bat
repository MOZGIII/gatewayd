@echo off

set dir=%~dp0
for %%* in (%dir%.) do set project=%%~n*

rem Unix mkdir compat
mkdir %dir%.goworkspace
mkdir %dir%.goworkspace\src
mklink /D %dir%.goworkspace\src\%project% ..\..
