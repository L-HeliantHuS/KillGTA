# KillGTA (一个用来打GTAOL犯罪之神挑战的小工具)
GTAOL, 犯罪之神挑战

抢劫首脑: 玩家所在的队伍成员依次完成五个难度为困难的抢劫, 顺序为: `全福` -> `越狱` -> `人道` -> `首轮` -> `太平洋`, 即可达成条件

末日首脑: 末日首脑分2个人的(犯罪之神2), 3个人的(犯罪之神3), 4个人的(犯罪之神4), 同一小队按顺序依次完成前置任务(可死亡),  准备任务(不可死亡)以及分红任务(不可死亡), 难度全为困难, 即可达成条件.

### 目的
由于硬性的条件是**不可死亡**, 但GTAOL当在出现死亡的一瞬间, 结束掉游戏进程, 不让它上传死亡数据, 即可重新上线继续从当前进度开始打, 不过要队伍里面的人都结束掉进程, 如果有人没结束掉, 那么他的进度将会重置(一个人重置之后所有人都要重新打), 为了防止这种尴尬的情况, 就有了此软件.

### 使用方法
程序使用`F4`键, `F4`键, `F4`键来干掉GTA5进程.


##### 修改默认快捷键
只需要启动的时候加上`-key`并指定按键就可以, 列如: `Client.exe -key=f3`, 即可修改快捷键为`f3`.

##### Online
启动Server, 默认会监听`25155`端口, 要将这个端口映射到公网, 起码客户端要能连接到你这个端口.

启动Client, 输入Server的地址, 列如服务端的ip为`192.168.5.2`, 监听的端口为`25155`, 则Client输入Server地址的时候就是`192.168.5.2:25155`, 连接成功后, Server会收到消息并提示, Client也会提示连接成功

##### Offline
如果没公网ip,或者和队友都有足够快的手速, 可以选择以离线模式运行本程序, 只需要在`Client.exe`同级目录打开`cmd`, 输入`Client.exe -online=false`, 即可直接运行本程序.

### 功能
`Server.go`
```
服务端, 用于分发任意一个客户端提交上来的kill请求, 并且支持心跳检测, 以及客户端数量检测.
```

`Client.go`

```
客户端, 使用TCP连接至服务端, 并每隔30秒都进行心跳检测, 当按下F4键, 向服务端发送一个kill请求, 服务器收到kill请求, 就会分发至各个客户端, 保证每个人都可以结束掉GTA.exe
```

### 特别感谢
[JetBrains](http://jetbrains.com/)提供IDE, 以用来编写此软件代码.

