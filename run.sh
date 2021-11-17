docker image rm lyt1123/open_im_server
tar zxvf Open-IM-Server.tgz -C ./
rm   -rf  Open-IM-Server.tgz
docker-compose down && docker-compose up -d --build