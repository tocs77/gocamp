@echo off
echo Generating TLS certificate with SANs...
openssl req -x509 -newkey rsa:2048 -nodes ^
  -keyout key.pem ^
  -out cert.pem ^
  -days 365 ^
  -subj "/CN=sch-grpc" ^
  -addext "subjectAltName=DNS:sch-grpc,DNS:localhost,IP:127.0.0.1"
echo.
echo Certificate files generated:
echo   - key.pem (private key)
echo   - cert.pem (certificate with SAN: sch-grpc, localhost, 127.0.0.1)
