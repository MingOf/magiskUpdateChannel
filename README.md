# Magisk

Magisk 是一套开放源代码的 Android 自定义工具套组，内置了 Magisk Manager、Root、启动脚本、SElinux 补丁和启动时认证 /dm-verity/ 强制加密移除等功能。


# 搭建 Magisk beta 更新通道

1.下载源文件

```bash
git clone https://github.com/MingOf/magiskUpdateChannel.git
cd magiskUpdateChannel
chmod a+x magiskChannel
```


2.启动
```bash
./magiskChannel [-d=域名或者主机IP] [-p=端口] [--listenAddress=0.0.0.0] [--listenPort=监听端口]
```
域名或者主机 IP 决定生成的配置的访问域名和端口

其中监听端口和监听地址是服务监听的地址和端口，默认监听 127.0.0.1:8080。
监听端口默认和配置文件访问端口[-p]相同，也可以分开指定，以便使用 nginx 进行端口转发 

指定域名后确保域名已经解析到自己的主机

启动后会定时拉取 magisk.apk 和 magisk.zip 到本地，保证最新版本

3.访问

```
# 配置文件
http://yourdomain/beta.json

# 网页
http://yourdomain/

# magisk.apk
http://yourdomain/magisk.apk

# magisk.zip
http://yourdomain/magisk.zip
```