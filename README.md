# Team2048_Tiktok
字节跳动青训营结营项目——极简抖音

本项目是基于HTTP框架Gin、ORM框架GORM的极简版抖音服务端项目

# 技术选型
- 项目代码采用**handler层**、**logic层**、**dao层**三层结构
- 使用**Snowflake**分布式ID生成器，它可以生成全局唯一、有序、可排序的ID
- 使用**GORM**操作MySQL数据区，并且优化GORM操作达到**防止SQL注入**效果
- 使用**MySQL**数据库对主要数据进行持久化，并为每个数据库表建立了**索引**
- 使用**Redis**数据库对热点数据进行持久化，使用**事务**保证数据一致性，提升服务性能
- 使用**JWT**鉴权，**MD5**密码加密，**ffmpeg**获取视频第一帧当作视频封面
- 使用**Viper**库进行配置文件的加载和配置管理
- 使用**Logger日志器**辅助项目的开发和维护
- 进行了**单元测试**，并且项目内置生成项目的**Swagger接口文档**


# 项目架构

![tiktok架构](https://user-images.githubusercontent.com/114276877/221083731-158a8876-d2c6-4d7e-af2a-1faa87921296.png)
将项目的处理逻辑分散到每一层，明确每一层的职责：

- **handler**负责参数校验和路由转发
- **logic**负责代码逻辑的实现
- **dao**负责数据的持久化

降低函数之间的耦合程度，使得项目呈现出“**低耦合，高内聚**”的特点


# 功能介绍

本项目主要实现以下功能：

视频：视频推送、视频投稿、发布列表

用户：用户注册、用户登录、用户信息

点赞：点赞操作、点赞列表

评论：评论操作、评论列表

关注：关注操作、关注列表、粉丝列表



# 数据库表设计

![Tiktok数据库ER图](https://user-images.githubusercontent.com/114276877/221083992-43e518c1-c5d1-4821-b42f-e1a6937f0e33.png)
![image](https://user-images.githubusercontent.com/114276877/221084025-92101067-f6b3-4a46-b238-12ba67c2a7b3.png)

# 项目启动

下载代码：git clone [https://github.com/FrancisChoi02/Team2048_Tiktok/.git](https://github.com/FrancisChoi02/Team2048_Tiktok/.git)

进入项目目录：cd your_repository

安装依赖：go mod init

配置环境变量：将MySQL数据库的地址、用户名和密码配置到环境变量中

启动项目：go run ./main.go

# Swagger接口文档启动

本项目使用Swagger生成接口文档，启动步骤如下：

安装Swagger：go get -u github.com/swaggo/swag/cmd/swag

为接口代码添加注释

在gin框架路由中进行注册：r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

生成Swagger文档: swag init

访问文档：运行项目，在浏览器打开[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
