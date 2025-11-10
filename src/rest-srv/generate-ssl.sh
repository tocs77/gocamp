#!/bin/bash

echo "Generating SSL certificate..."
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost"

echo ""
echo "Certificate files generated:"
echo "  - key.pem (private key)"
echo "  - cert.pem (certificate)"

