@echo off
echo Testing HTTPS endpoint...
curl.exe -k -v https://localhost:8000/orders
pause



