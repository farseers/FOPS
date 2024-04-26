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

docker run --name fops  -d --network=host -e "Database_default=DataType=mysql,PoolMaxSize=50,PoolMinSize=1,ConnectionString=root:harlen@tcp(host.docker.internal:3307)/fops2?charset=utf8&parseTime=True&loc=Local" -e "Redis_default=Server=host.docker.internal:6379,DB=13,Password=,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000" -e "FOPS_ConnString=DataType=clickhouse,PoolMaxSize=50,PoolMinSize=1,ConnectionString=clickhouse://default:@host.docker.internal:9000/linkTrace?dial_timeout=10s&read_timeout=20s" hub.fsgit.cc/fops:483