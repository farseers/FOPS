# 注意，这里的构建上下文，是在git源代码的根目录
FROM golang:1.20.11-alpine AS build
# 设置github代理
ENV GOPROXY https://goproxy.cn,direct
# 进入到项目目录中
WORKDIR /src/fops
# 复制go.mod文件
COPY ./fops/go.mod .
# 下载依赖（支持docker缓存）
RUN go mod download
# 将源代码复制到此
COPY ./fops .
# 删除go.work文件
#RUN rm -rf go.work
# 更新go.sum
RUN go mod tidy
# farseer项目
WORKDIR /src/farseer-go
COPY ./farseer-go .
# 进入到项目目录中
WORKDIR /src/fops
# 编译
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/fops-server -ldflags="-w -s" .

FROM steden88/cicd:1.0 AS base
WORKDIR /app
COPY --from=build /app .
# 复制配置（没有配置需要注释掉）
COPY --from=build /src/fops/farseer.yaml .
# 复制视图（没有视图需要注释掉）
#COPY --from=build /src/views ./views
# 复制静态资源（没有静态资源需要注释掉）
#COPY --from=build /src/fops/wwwroot ./wwwroot
# 复制vue源码
COPY --from=build /src/fops/wwwsrc ./wwwsrc
WORKDIR /app/wwwsrc
# 构建npm
RUN npm install
RUN npm run build

# 创建目录
RUN mkdir -p /app/wwwroot/
# 前端文件移到静态目录
RUN cp -r /app/wwwsrc/dist/* /app/wwwroot/
# 删除源文件
RUN rm -rf /app/wwwsrc

#设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai    /etc/localtime

WORKDIR /app
ENTRYPOINT ["./fops-server"]

