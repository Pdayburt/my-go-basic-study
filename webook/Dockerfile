#基础镜像
FROM ubuntu:20.04
#将编译后的webook二进制文件打包进此镜像
COPY webook /app/webook
#将/app设为工作目录
WORKDIR /app
CMD ["/app/webook"]