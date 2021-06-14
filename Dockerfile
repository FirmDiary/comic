#使用镜像大小只有5m的alpine
FROM alpine:latest

#设置工作路径
WORKDIR /go/src

#把上述编译好的main文件添加到镜像里
COPY . .

#暴露容器内部端口
EXPOSE 9090

CMD ["sh", "/www/comic/backend/script/build.sh"]