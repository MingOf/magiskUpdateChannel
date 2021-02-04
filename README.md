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
./magiskChannel [-d=域名或者主机IP] [-p=端口] [--listenAddress=127.0.0.1] [--listenPort=监听端口]
```
-d: 域名或者主机 IP，决定生成的配置的访问域名，默认使用主机的公网 IP 地址 

-p: 端口，决定生成的配置的访问域名端口，默认 80

如：
```bash
./magiskChannel -d example.com -p 80 

=>

http://exmaple.com:80/beta.json

http://example.com:80/

```

--listenAddress: 启动 http 服务的监听地址，默认 0.0.0.0

--listenPort: 启动 http 服务的监听端口，默认 80

通常情况下不用特别指定 listenAddress 和 listenPort

如果使用 nginx 可以将 listenAddress 改为 127.0.0.1，listenPort 改为其他端口

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