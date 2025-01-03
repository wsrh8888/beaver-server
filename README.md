
## 项目介绍
+ 本项目使用go-zero做的一款微服务项目，仓库代码完全是基于商业化代码开发的。
+ 这是项目的第一个版本历时2个月， 前端使用uniapp， 后端使用go-zero
+ 加入QQ群：[![加入QQ群](https://img.shields.io/badge/加入QQ群-1013328597-blue.svg)](https://qm.qq.com/q/82rbf7QBzO)（1013328597）

## 服务端口
| 端口 | api | rpc | admin |
|:---------:|:--------:|:--------:|:--------:|
|user|21000|22000|23000|
|auth|21100|22100|23100|
|friend|21200|22200|23200|
|chat|21300|22300|23300|
|ws|21400|22400|23400|
|group|21500|22500|23500|
|file|21600|22600|23600|

### 功能介绍
+ 用户注册
+ 用户登录
+ 好友列表
+ 好友详情
+ 个人信息
+ 修改头像(七牛云)
+ 个人二维码
+ 扫码添好友
+ 私聊(文本、图片、表情)
+ ws消息
+ 删除好友
+ 网关

### 启动命令教程：

+ 安装etcd、mysql、redis
build/docker-compose.yaml
+ 初始化依赖
go mod tidy
+ 初始化数据库
go run main.go  -db
+ 本地运行
需要先启动RPC服务后在启动API服务

<img src="./static/1.png"/>
<img src="./static/2.png"/>




### 项目列表
| [GitHub仓库]    |  说明                                                                                      
| ------------------------------------------------------------ | --------------------------------------------------------------------------
| [beaver-server](https://github.com/wsrh8888/beaver-server)               | 后端服务  |
| [beaver-mobile](https://github.com/wsrh8888/beaver-mobile)        | uniapp客户端 |




### 应用截图

登录界面
<img src="./static/login.png"/>

注册界面界面
<img src="./static/register.png"/>

我的界面
<img src="./static/mine.png"/>

我的二维码
<img src="./static/qcode.png"/>

好友列表
<img src="./static/friend.png"/>

消息页面
<img src="./static/message.png"/>

聊天页面
<img src="./static/chat.png"/>

聊天页面
<img src="./static/chat1.png"/>

详情页面
<img src="./static/info.png"/>







