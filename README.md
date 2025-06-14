# 海狸IM


### 更新日期
2025-06-03

### 文档地址
https://wsrh8888.github.io/beaver-docs/

### 项目介绍
+ 本项目是利用工作之余做的一个聊天IM。
+  前端使用uniapp， 后端使用go-zero， 桌面端使用electron
+ 如果喜欢可以点一个star
+ 加入QQ群：[![加入QQ群](https://img.shields.io/badge/加入QQ群-1013328597-blue.svg)](https://qm.qq.com/q/82rbf7QBzO)（1013328597）



# 服务端口
| 端口 | api | rpc | admin |
|:---------:|:--------:|:--------:|:--------:|
|user|21000|22000|23000|
|auth|21010|22010|23010|
|friend|21020|22020|23020|
|chat|21030|22030|23030|
|ws|21040|22040|23040|
|group|21050|22050|23050|
|file|21060|22060|23060|
|emoji|21070|22070|23070|
|gateway|21080|-----|23080|
|moment|21090|-----|23090|
|system|21100|-----|23100|
|call|21110|-----|23110|
|settings|21120|-----|23120|
|notification|21140|-----|23140|
|feedback|21150|-----|23150|




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
| [GitHub仓库]    |   [Gitee仓库]    |说明                                                                                      
| ------------------------------------------------------------ | --------------------------------------------------------------------------|--------------------------------------------------------------------------|
| [beaver-server](https://github.com/wsrh8888/beaver-server)               |[beaver-server](https://gitee.com/dawwdadfrf/beaver-server)               | 后端服务  |
| [beaver-mobile](https://github.com/wsrh8888/beaver-mobile)        | [beaver-mobile](https://gitee.com/dawwdadfrf/beaver-mobile)               |手机端 |
| [beaver-desktop](https://github.com/wsrh8888/beaver-desktop)        | [beaver-desktop](https://gitee.com/dawwdadfrf/beaver-desktop)               |桌面端 |


### 更新记录
20240426
- 1、增加通过会话id获取会话信息接口
- 2、ws代理接口增加会话id参数
- 3、群组详细接口增加组人数
- 4、部分ws消息转发时候增加会话id
- 5、修复部分字段参数异常
- 6、redis去掉密码校验

20250422
- 1、增加反馈功能
- 2、增加群聊功能
- 3、优化各种bug
- 4、头像变更为id服务端做转发
- 5、go-zero版本升级
- 6、好友模块优化
- 7、最近会话列表优化

20250119
- 1、修改部分属性字段比如Id变更为ID
- 2、修改user表中的user_id 变更为uuid
- 3、增加朋友圈表以及相关服务
- 4、增加表情包表以及相关服务
- 5、ws服务重构，拆分为不同的模块
- 6、增加group表以及相关服务


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
<img src="./static/home.png"/>

聊天页面
<img src="./static/chat.png"/>

群聊页面
<img src="./static/groupChat.png"/>

群聊详情
<img src="./static/groupConfig.png"/>

群聊移除界面
<img src="./static/removeMember.png"/>

详情页面
<img src="./static/info.png"/>

设置页面
<img src="./static/setting.png"/>


关于页面
<img src="./static/about.png"/>

问题反馈
<img src="./static/feedback.png"/>
