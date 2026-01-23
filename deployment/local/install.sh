#!/bin/bash

set -eo pipefail

echo "---------------------------------------------------------------"
echo "INSTALL LocalPaaS LOCALLY"
echo "---------------------------------------------------------------"

# Delete all unused data that take the disk space
# docker system prune -a -f

LOCALPAAS_DIR=.appdata/localpaas
LOCALPAAS_CERTS=$LOCALPAAS_DIR/certs

mkdir -p $LOCALPAAS_DIR
mkdir -p $LOCALPAAS_CERTS

NGINX_ETC=$LOCALPAAS_DIR/nginx/etc
NGINX_LOG=$LOCALPAAS_DIR/nginx/log
NGINX_SHARE=$LOCALPAAS_DIR/nginx/share
NGINX_CERTS=$NGINX_SHARE/certs

mkdir -p $NGINX_ETC/conf.d
mkdir -p $NGINX_LOG
mkdir -p $NGINX_SHARE
mkdir -p $NGINX_SHARE/default
mkdir -p $NGINX_SHARE/domains
mkdir -p $NGINX_SHARE/html
mkdir -p $NGINX_CERTS/fake

# Remove all app config files from nginx
echo "Remove all app config files from nginx..."
rm -rf $NGINX_ETC/conf.d/*.conf

# Download nginx conf files
echo "Copy nginx config files..."
cp deployment/local/nginx/nginx.conf $NGINX_ETC/nginx.conf
cp deployment/local/nginx/localpaas.conf $NGINX_ETC/conf.d/localpaas.conf

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

# Init docker swarm
echo "Init docker swarm..."
docker swarm init || true

# Create overlay network for nginx to discover services
echo "Create overlay network 'localpaas_net'..."
docker network create --driver overlay --attachable localpaas_net || true

# Deploy localpaas stack
echo "Deploy localpaas stack..."
cp deployment/local/app_stack_nginx.yaml $LOCALPAAS_DIR/../localpaas.yaml
docker stack deploy -c $LOCALPAAS_DIR/../localpaas.yaml localpaas

sleep 5
make seed-data-with-clear

echo "---------------------------------------------------------------"
echo "DONE."
echo "---------------------------------------------------------------"
