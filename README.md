## stool 使用说明

1、创建目录
```shell
mkdir -p /data/stool
```

2、从 release 下载 stool 执行文件
```shell
wget https://github.com/LeeEirc/stool/releases/download/v0.0.4/stool-linux-amd64 -O /data/stool/stool-linux-amd64
```

3、创建配置文件参考[config.yml](config.yml)
```shell
touch /data/stool/stool.yaml
```

4、启动服务
```shell
/data/stool/stool-linux-amd64 -c /data/stool/stool.yaml
```
