#!/bin/bash

# Step 1: Get nginx config file from nginx container
CONFIG_DIR="/etc/nginx/conf.d"
TMP_FILE="nginx-temp.conf"
CONTAINER_NAME="nginx"

# Check if the branch name is "master"
#if [ "$CI_COMMIT_REF_NAME" = "master" ]; then
#    FILEPATH="$CONFIG_DIR/$CI_PROJECT_NAME.0xnina.xyz.conf"
#    SERVER_NAME="$CI_PROJECT_NAME.0xnina.xyz"
#else
FILEPATH="$CONFIG_DIR/$APP_NAME.fidesy.xyz.conf"
SERVER_NAME="$APP_NAME.fidesy.xyz"
#fi

CONTENT="""
server {
    listen 80;
    listen [::]:80;

    server_name $SERVER_NAME;

    location / {
      proxy_set_header    Host                \$http_host;
      proxy_set_header    X-Real-IP           \$remote_addr;
      proxy_set_header    X-Forwarded-For     \$proxy_add_x_forwarded_for;


      proxy_pass http://$SERVER_NAME:${SWAGGER_PORT};
    }
}
"""
echo $CONTENT > ${TMP_FILE}

# Copy the updated configuration back to the nginx container
docker cp ${TMP_FILE} ${CONTAINER_NAME}:${FILEPATH}

# Reload NGINX within the container to apply the changes
docker exec ${CONTAINER_NAME} nginx -s reload

# Remove TMP file
rm ${TMP_FILE}
