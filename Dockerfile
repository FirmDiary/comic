#使用镜像大小只有5m的alpine
FROM alpine:latest

#创建环境变量
RUN mkdir -p /go/app

#设置工作路径
WORKDIR /go/app

#把上述编译好的main文件添加到镜像里
COPY . .

#暴露容器内部端口
EXPOSE 9090

#入口
ENTRYPOINT ["/go/app/aisvc"]