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

### 使用
go build -o integrationService

./integrationService console filesRemove 202101 202202

