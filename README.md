# au

Analysis url form log file.

## 功能

- 生成域名url列表
- 统计组件使用次数
- 使用组件的url列表

## 使用方法
```
// 使用help获取帮助信息
./au --help
```
## 示例

```
// type 参数必选 1：生成域名url列表    2：统计组件使用次数    3：使用组件的url列表
// pattern 参数可选  后带字符串可以是正则，一般用来匹配日志文件夹中某些日志，比如某一天的日志

./au --path path/to/log --type 1 --pattern "xxx_processor.log.2017-09"

```