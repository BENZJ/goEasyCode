# goEasyCode
## 简介
go实现命令行读取数据库，生成相关代码
## 项目结构
## 版本
### v0.0.1
  - 读取数据库表结构
- 使用模板引擎自动生成

## 参考项目
- [converter](https://github.com/gohouse/converter) go项目，读取数据库生成golang entity文件
- [EasyCode插件](https://github.com/makejavas/EasyCode)


## 使用说明

```bash 
## 编译
go build -o easycode.exe  cmd/cli.go 

##使用dsn配置
easycode --dsn esop:EsopManDev123\!@tcp\(172.15.6.184:3306\)/wealth_service?charset=utf8 --table org_service_log_coordination -template Do

## yaml配置读取 yaml配置模板参考mysql-templ.yaml
easycode --table org_service_log_coordination -template Do
```