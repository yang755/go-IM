# go-IM

go语言 IM聊天工具

1、使用 下载代码

2、cd ./IM/; go build -o server main.go  server.go  user.go; 此命令编译成server可执行文件

3、输入./server 执行即可，此时server已启动

4、客户端连接可以使用nc 127.0.0.1 8888或者使用telnet ip 8888连接



```shell
[xxxMacBook-Air:~ zhangyang$ nc 127.0.0.1 8888
[127.0.0.1:2058]127.0.0.1:2058:已上线
who
[127.0.0.1:2058]127.0.0.1:2058:在线...
[127.0.0.1:50858]127.0.0.1:50858:在线...
[127.0.0.1:50866]127.0.0.1:50866:在线...
[127.0.0.1:50876]127.0.0.1:50876:在线...
```

