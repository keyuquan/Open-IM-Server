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

echo "压缩文件"
tar -zcvf .build/Open-IM-Server.tgz bin/* script/*  ./Dockerfile  ./docker-compose.yml config/*  cmd/*
#rm registerserver
echo "scp 文件到服务器"
scp .build/Open-IM-Server.tgz aliyun-stone:/root/Open-IM-Server

echo "完事"

chmod +x ./run.sh
echo "scp run.sh 文件到服务器"
scp ./run.sh aliyun-stone:/root/Open-IM-Server
