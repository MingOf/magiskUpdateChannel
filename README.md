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
./magiskChannel [-d=域名或者主机IP] [-p=端口]
```
域名或者主机 IP 会以及端口会传递给 index.tmpl 模板

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