@echo off
cd /d "%~dp0"
if exist proto\gen rmdir /s /q proto\gen
mkdir proto\gen
docker compose -f compose-grpc-build.yml run --rm protoc --build
