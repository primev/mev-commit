#!/usr/bin/env sh

# This script automates the creation of a Certificate Authority (CA) and a server
# certificate. It generates a private key and self-signed certificate for the CA,
# then creates a server private key and CSR, and signs the CSR with the CA's
# private key to produce a server certificate. The script requires OpenSSL and optionally
# accepts parameters for the CA subject, server subject, and a server extension config file.
# If the server extension config file is not provided, a default file (server-ext.cnf)
# will be generated with basic constraints and key usage. For more information about the
# configuration file, see: https://www.openssl.org/docs/manmaster/man5/x509v3_config.html

# Parameters:
# 1. CA Subject (optional): The subject details for the CA certificate. Default: "/C=US/O=Server's CA"
# 2. Server Subject (optional): The subject details for the server certificate. Default: "/C=US/O=Server"
# 3. Server Extension Config File (optional): Path to a configuration file with extensions
#    for the server certificate. If not provided, a default configuration is used.

# Generated files:
# - CA private key (ca-key.pem)
# - CA self-signed certificate (ca-cert.pem)
# - Server private key (server-key.pem)
# - Server CSR (server-req.pem)
# - Server certificate (server-cert.pem)

# Usage:
# Execute this script with up to three optional arguments:
# ./script.sh [CA Subject] [Server Subject] [Server Extension Config File]
# Example: ./script.sh "/C=US/O=My CA" "/C=US/CN=myserver.example.com" "myserver-ext.cnf"

# Ensure OpenSSL is installed and accessible, prepare the optional server-ext.cnf file
# with necessary server certificate extensions if desired, and execute this script.
# Verify output for CA and server certificate details.

# Note: Designed for educational or development purposes. Adapt carefully for production use.


CA_KEY="ca-key.pem"
CA_CERT="ca-cert.pem"
SERVER_KEY="server-key.pem"
SERVER_REQ="server-req.pem"
SERVER_CERT="server-cert.pem"
CA_SUBJ=${1:-"/C=US/O=Server's CA"}
SERVER_SUBJ=${2:-"/C=US/O=Server"}
SERVER_EXT=${3:-"server-ext.cnf"}

# Generate a default server-ext.cnf file if not provided.
if [ ! -f "${SERVER_EXT}" ]; then
  echo "No server extension conf provided; generating a default configuration:"
cat << EOH > "${SERVER_EXT}"
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment
EOH
  cat "${SERVER_EXT}"
fi

# Generate CA's private key and self-signed certificate.
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout "${CA_KEY}" -out "${CA_CERT}" -subj "${CA_SUBJ}"
echo "CA's self-signed certificate:"
openssl x509 -in "${CA_CERT}" -noout -text

# Generate server's private key and certificate request (CSR).
openssl req -newkey rsa:4096 -nodes -keyout "${SERVER_KEY}" -out "${SERVER_REQ}" -subj "${SERVER_SUBJ}"
# Use CA's private key to sign server's CSR and generate the server's certificate.
openssl x509 -req -in "${SERVER_REQ}" -days 365 -CA "${CA_CERT}" -CAkey "${CA_KEY}" -CAcreateserial -out "${SERVER_CERT}" -extfile "${SERVER_EXT}"
echo "Server's CA signed certificate:"
openssl x509 -in "${SERVER_CERT}" -noout -text
