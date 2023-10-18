# 权限监听

在一些业务中，我们会将日志输出到文件，但在 ThinkPHP 等一些框架上的使用，会导致文件权限的混乱，这里我们采用实时监听的方式，当出现异常时，自动修复权限，目前仅支持
Linux 系统。

```shell
# 安装
$ curl -Ls https://github.com/hongfs/perm-listen/releases/download/v0.0.3/listen -o /usr/local/bin/listen
$ chmod +x /usr/local/bin/listen
```

运行需要使用的环境变量：

```shell
# 监听的目录
$ export LISTEN_PATH=/data/wwwroot
# 监听的用户，给这个B权限
$ export LISTEN_USER=www
# 监听的文件扩展
$ export LISTEN_EXTENSION=.log
```

后台运行

```shell
$ LISTEN_PATH=/data/wwwroot LISTEN_EXTENSION=.log LISTEN_USER=www /usr/local/bin/listen &
```

查看日志

日志文件位于： `/tmp/perm-listen.log`，可以使用 `tail -f /tmp/perm-listen.log` 实时查看。

杀死进程

```shell
$ ps aux | grep /usr/local/bin/listen
root     20705  0.0  0.0 1227980 3312 pts/0    Sl   20:21   0:00 /usr/local/bin/listen

# 干他
$ kill -9 20705
```
