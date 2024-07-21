# 更新farseer-go框架
#cd ../farseer-go && sh git-update.sh
# 更新fops
#cd ../fops && git pull
# 将忽略文件复制到上下文根目录中
#\cp .dockerignore ../
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./fops -ldflags="-w -s" .
# 编译
docker build -t hub.fsgit.cc/hub:fops.672 --network=host .
# 发到内网
docker push hub.fsgit.cc/hub:fops.672 && docker rmi hub.fsgit.cc/hub:fops.672

# docker
docker service rm fops
docker service create --name fops --replicas 1 -d --network=net \
 --with-registry-auth \
--constraint node.role==manager \
--mount type=bind,src=/etc/localtime,dst=/etc/localtime \
--mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
--mount type=bind,src=/home/nfs/fops,dst=/var/lib/fops \
-l "traefik.http.routers.fops.rule=Host(\`fops.fsgit.cc\`)" \
-l "traefik.http.routers.fops.entrypoints=websecure" \
-l "traefik.http.routers.fops.tls=true" \
-l "traefik.http.services.fops.loadbalancer.server.port=8889" \
hub.fsgit.cc/fops:dev

docker run --name fops  -d --network=host -e "Database_default=DataType=mysql,PoolMaxSize=50,PoolMinSize=1,ConnectionString=root:harlen@tcp(host.docker.internal:3307)/fops2?charset=utf8&parseTime=True&loc=Local" -e "FOPS_ConnString=DataType=clickhouse,PoolMaxSize=50,PoolMinSize=1,ConnectionString=clickhouse://default:@host.docker.internal:9000/linkTrace?dial_timeout=10s&read_timeout=20s" hub.fsgit.cc/fops:483