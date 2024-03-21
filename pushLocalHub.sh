# 更新farseer-go框架
cd ../farseer-go && sh git-update.sh
# 更新fops
cd ../fops && git pull
# 将忽略文件复制到上下文根目录中
#\cp .dockerignore ../
# 编译
docker build -t hub.fsgit.cc/fops:dev --network=host -f ./Dockerfile ../
# 发到内网
docker push hub.fsgit.cc/fops:dev && docker rmi hub.fsgit.cc/fops:dev

# docker
docker service rm fops
docker service create --name fops --replicas 1 -d --network=net \
--constraint node.role==manager \
--mount type=bind,src=/etc/localtime,dst=/etc/localtime \
--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
--mount type=bind,src=/home/nfs/fops,dst=/var/lib/fops \
-l "traefik.http.routers.fops.rule=Host(\`fops.fsgit.cc\`)" \
-l "traefik.http.routers.fops.entrypoints=websecure" \
-l "traefik.http.routers.fops.tls=true" \
-l "traefik.http.services.fops.loadbalancer.server.port=8889" \
steden88/fops:2.0.beta3