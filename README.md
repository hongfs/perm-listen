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
$ LISTEN_PATH=/data/wwwroot
# 监听的用户，给这个B权限
$ LISTEN_USER=www
# 监听的文件扩展
$ LISTEN_EXTENSION=.log
```

运行

```shell
$ LISTEN_PATH=/data/wwwroot LISTEN_EXTENSION=.log LISTEN_USER=www /usr/local/bin/listen
```

日志文件位于： `/tmp/perm-listen.log`，可以使用 `tail -f /tmp/perm-listen.log` 实时查看。

