#!/bin/bash

set -e

VERSION=$1
CUR_VERSION=""
IMAGE_NAME=open-im-server
if test -z $VERSION; then
   echo "请输入版本号!!!"
   exit 1
fi

if [ "$(echo $VERSION | grep "beta-")" != "" ]; then
  CUR_VERSION="develop"
fi

if [ "$(echo $VERSION | grep "v-")" != "" ]; then
  CUR_VERSION="release"
fi

if [ "$CUR_VERSION" = "" ]; then
  echo "版本号错误,无法推送!!!"
  exit 1
fi

echo "开始编译..."
cd script/
./build_all_service.sh
echo "编译完成..."

cd ..

# 写入版本号
echo "写入版本号"
echo "$1" > ./config/VERSION

echo "生成Dockerfile文件"
cat>.build/Dockerfile<<EOF
FROM ubuntu

RUN rm -rf /var/lib/apt/lists/*
RUN apt-get update && apt-get install apt-transport-https && apt-get install procps\
&&apt-get install net-tools
#Non-interactive operation
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y vim curl tzdata gawk
#Time zone adjusted to East eighth District
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata


#set directory to map logs,config file,script file.
VOLUME ["/Open-IM-Server/logs","/Open-IM-Server/config","/Open-IM-Server/script","/Open-IM-Server/db/sdk"]

WORKDIR /Open-IM-Server/script

CMD ["./docker_start_all.sh"]
EOF

echo "生成docker-compose.yml文件"
cat>.build/docker-compose.yml<<EOF
version: '3'
services:
  open-im-server:
    image: open_im_server
    build: .
#      context: .
#      dockerfile: deploy.Dockerfile
    container_name: open-im-server
    volumes:
      - ./logs:/Open-IM-Server/logs
      - ./config/config.yaml:/Open-IM-Server/config/config.yaml
      - ./db/sdk:/Open-IM-Server/db/sdk
      - ./script:/Open-IM-Server/script
      - ./bin:/Open-IM-Server/bin
    restart: always
    environment:
      - CONFIG_FILE=../config/config.yaml
    networks:
      - traefik
    labels:
      - "traefik.enable=true"
      - "traefik.docker.network=traefik"
      # Entry Point for https
      - "traefik.http.routers.im-api.entrypoints=http"
      - "traefik.http.routers.im-api.rule=Host(\`appstone.top\`) && PathPrefix(\`/imapi\`)"
      - "traefik.http.routers.im-api.middlewares=im-api-stripPrefix@file"
      - "traefik.http.routers.im-api.service=im-api-service"
      - "traefik.http.services.im-api-service.loadbalancer.server.port=10000"
      # websocket
      - "traefik.http.routers.im-api-ws.entrypoints=http"
      - "traefik.http.routers.im-api-ws.rule=Host(\`appstone.top\`) && PathPrefix(\`/imapiws\`)"
      - "traefik.http.routers.im-api-ws.middlewares=im-api-ws-stripPrefix@file"
      - "traefik.http.routers.im-api-ws.service=im-api-ws-websocket"
      - "traefik.http.services.im-api-ws-websocket.loadbalancer.server.port=17778"
    logging:
      driver: json-file
      options:
        max-size: "1g"
        max-file: "2"

networks:
  traefik:
    external: true
EOF

echo "压缩文件"
tar -zcvf .build/all.tgz bin/* script/*  db/* .build/Dockerfile .build/docker-compose.yml config/*
#rm registerserver
echo "scp 文件到服务器"
scp .build/all.tgz aliyun-stone:/root/Open-IM-Server

# run.sh
echo "生成 run.sh"
cat>.build/run.sh<<EOF
tar zxvf all.tgz -C ./
mv .build/* ./
docker-compose down && docker-compose up -d --build  && rm all.tgz
EOF

#echo "更新最新版本库"
#docker tag $(docker images | grep ${IMAGE_NAME} | head -1 | awk '{print $3}') ${IMAGE_NAME}:$CUR_VERSION
#
#echo "清理本地无用镜像..."
#docker rmi -f $(docker images | grep ${IMAGE_NAME} | awk '{print $3}')

echo "完事"

chmod +x .build/run.sh
echo "scp run.sh 文件到服务器"
scp .build/run.sh aliyun-stone:/root/Open-IM-Server
