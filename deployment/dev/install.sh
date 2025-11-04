#!/bin/bash

set -eo pipefail

echo "---------------------------------------------------------------"
echo "INSTALL LocalPaaS"
echo "---------------------------------------------------------------"

NGINX_ETC=nginx/etc
NGINX_LOG=nginx/log
NGINX_SHARE=nginx/share
NGINX_CERTS=$NGINX_SHARE/certs

mkdir -p $NGINX_ETC/conf.d
mkdir -p $NGINX_LOG
mkdir -p $NGINX_SHARE
mkdir -p $NGINX_SHARE/default
mkdir -p $NGINX_SHARE/domains
mkdir -p $NGINX_CERTS/fake

# Download nginx conf files
echo "Download nginx config files..."
curl -sL "https://raw.githubusercontent.com/localpaas/localpaas/dev/deployment/dev/nginx/nginx.conf" -o $NGINX_ETC/nginx.conf
curl -sL "https://raw.githubusercontent.com/localpaas/localpaas/dev/deployment/dev/nginx/localpaas.conf" -o $NGINX_ETC/conf.d/localpaas.conf

# Gen self-signed SSL certs
if [ ! -f "$NGINX_CERTS/fake/local.key" ]; then
  echo "File 'fake/local.key' does not exist. Generate new file..."
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout $NGINX_CERTS/fake/local.key -out $NGINX_CERTS/fake/local.crt \
    -subj "/CN=*.swarm.localhost"
fi

# Gen dhparam.pem
if [ ! -f "$NGINX_CERTS/dhparam.pem" ]; then
  echo "File 'dhparam.pem' does not exist. Generate new file..."
  openssl dhparam -out $NGINX_CERTS/dhparam.pem 2048
fi

# Create overlay network for nginx to discover services
echo "Create overlay network 'localpaas_net'..."
docker network create --driver overlay --attachable localpaas_net || true

# Download app_stack_nginx.yaml
echo "Download app_stack_nginx.yaml..."
curl -sL "https://raw.githubusercontent.com/localpaas/localpaas/dev/deployment/dev/app_stack_nginx.yaml" -o localpaas.yaml

# Deploy localpaas stack
echo "Deploy localpaas stack..."
docker pull localpaas/localpaas-dev:app-latest # pull latest image
docker stack deploy -c localpaas.yaml localpaas

sleep 5
docker run --net localpaas_internal_net \
  -e LP_PLATFORM=remote -e LP_DB_HOST=db -e LP_DB_PORT=5432 -e LP_DB_DB_NAME=localpaas \
  -e LP_DB_USER=localpaas -e LP_DB_PASSWORD=abc123 -e LP_DB_SSL_MODE=disable \
  -w /app localpaas/localpaas-dev:app-latest \
  make seed-data-with-clear

echo "---------------------------------------------------------------"
echo "DONE."
echo "---------------------------------------------------------------"
