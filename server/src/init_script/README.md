## 目录结构

├── node-offlinetx
│   └── signTx.js
├── otc
│   ├── conf
│   │    └── app.conf
│   └── otc
├── otc_admin
│   ├── conf
│   │   └── app.conf
│   └── otc_admin
├── README.md
└── sql
    ├── create_db.sql
    ├── otc
    │   ├── eusd
    │   │   └── __all_table_create.sql
    │   ├── game
    │   │   └── __all_table_create.sql
    │   ├── otc
    │   │   └── __all_table_create.sql
    │   └── usdt
    │       └── __all_table_create.sql
    ├── otc_admin
    │   └── __all_table_create.sql
    └── otc_admin_init.sql

## 1 应用部署

### 1.1 otc 配置

| 配置名称 | 配置说明 | 备注 |
| ---- | ---- | ---- |
| appurl  |  第三方应用跳转链接 |  appurl=https://47.75.178.99 |
| database::default | default数据库实例配置 |  修改为生产环境的default数据库实例配置 |
| database::otc | otc数据库实例配置 |  修改为生产环境的otc数据库实例配置 |
| database::otc_admin | otc_admin数据库实例配置 |  修改为生产环境的otc_admin数据库实例配置 |
| redis::host | redis服务地址 | 修改为生产环境实际配置 |
| redis | redis服务端口 | 修改为生产环境实际配置 |

### 1.2 otc_admin 配置

| 配置名称 | 配置说明 | 备注 |
| ---- | ---- | ---- |
| appurl  |  第三方应用跳转链接 |  appurl=https://47.75.178.99 |
| database::default | default数据库实例配置 |  修改为生产环境的default数据库实例配置 |
| database::otc | otc数据库实例配置 |  修改为生产环境的otc数据库实例配置 |
| database::otc_admin | otc_admin数据库实例配置 |  修改为生产环境的otc_admin数据库实例配置 |
| redis::host | redis服务地址 | 修改为生产环境实际配置 |
| redis | redis服务端口 | 修改为生产环境实际配置 |

## 2 数据库部署

### 创建数据库

>sql/create_db.sql

### 创建数据库表

#### otc_admin 实例

执行如下脚本：

>sql/otc_admin/__all_table_create.sql
>sql/otc_admin_init.sql

#### otc 实例

执行如下脚本：

>sql/otc/eusd/__all_table_create.sql
>sql/otc/game/__all_table_create.sql
>sql/otc/otc/__all_table_create.sql
>sql/otc/usdt/__all_table_create.sql

## 3 离线签名配置

查看node版本，如果未安装，参考 http://www.runoob.com/nodejs/nodejs-install-setup.html 进行安装

>node -v

安装bitcoin模块

>npm i bitcoinjs-lib
>npm init
>npm install express --save

启动服务

> node signTx.js

## 4 配置服务监控

上述 otc、otc_admin、node-offlinetx 服务都需配置服务监控，服务crash时能够自动拉起。

参考 supervisor 或其它。
