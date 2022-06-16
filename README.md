# 综合服务项目

>总体是为了做一些灵活性的功能，使用Go语言进行运行，编译运行速度快，支持并发执行，效率高

### 技术栈

1. golang 1.17
2. mysql
3. mongodb
4. redis等

### 功能
+ 接口api
+ 执行脚本

### 业务
+ 删除掉文件服务项目中多余的文件（删除jz生成的各种文件，支持时间范围删除功能）
+ 修复替换掉企微聊天数据中voice_meeting_call类型的部分错误数据
+ 修复company_user表格中的客户企微昵称为空的问题

### 使用
go build -o integrationService



### 功能

| 命令                                                       | 参数          | 备注                             |
| ---------------------------------------------------------- | ------------- | -------------------------------- |
| ./integrationService console filesRemove                   | 202101 202202 | 删除之前企微其他企业文件         |
| ./integrationService console repairMediaData msg_create_at | 2022-05-21    | 按日期替换 修复聊天语音          |
| ./integrationService console repairMediaData msg_id        | xxx           | 按消息id替换 修复聊天语音        |
| ./integrationService console repairWechatBind              | -             | 修复company_user企微绑定企微昵称 |
|                                                            |               |                                  |
|                                                            |               |                                  |
|                                                            |               |                                  |



