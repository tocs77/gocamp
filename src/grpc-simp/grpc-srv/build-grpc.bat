@echo off
cd /d "%~dp0"
docker compose -f compose-grpc-build.yml run --rm protoc --build
