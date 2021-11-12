tar zxvf Open-IM-Server.tgz -C ./
mv .build/* ./
docker-compose down && docker-compose up -d --build