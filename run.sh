docker image rm
tar zxvf Open-IM-Server.tgz -C ./
rm   -rf  Open-IM-Server.tgz
docker-compose down && docker-compose up -d --build