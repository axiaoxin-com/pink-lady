# webserver 包的使用方法

```
go get -u github.com/axiaoxin-com/pink-lady/src/webserver
```

大致流程如下，具体参考 [src/main.go](../main.go)

1. 加载配置文件，根据配置信息个性化 web server

```
webserver.InitWithConfigFile(path/to/config, configname, configtype)
```

2. 创建 app 路由

```
app := webserver.NewGinEngine(nil)
```

3. 启动 server

```
webserver.Run(app, yourAppURLRegister)
```
