@echo off
cd /d "%~dp0"
if exist protoc rmdir /s /q protoc
mkdir protoc
docker compose -f compose-grpc-build.yml run --rm protoc --build
