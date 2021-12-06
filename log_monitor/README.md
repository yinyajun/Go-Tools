# LogMonitor

日志监控服务，重构自项目https://github.com/qieangel2013/goMonitorLog

修改了算法实现和检测流程。

支持

* 实时关键词检测（使用AC自动机算法来支持多模式串文本匹配）

* 支持定制忽略词

* 支持分割词忽略（例如：匹配【机器学习】，分割词定义为“-”，那么日志中出现【机器-学习】，也能被成功检测）

* 支持带日期格式的日志文件

* 支持带限流器的钉钉报警器（令牌桶限流器）

  

### 配置说明

```toml
# 规则
rules = [
    { name = "rule1", ignored_words = [], error_words = ["abc"], warning_words = [] },
    { name = "rule2", ignored_words = [], error_words = ["123"], warning_words = [] },
]

# 监控日志文件
files = [
    { name = "/tmp/abc.log", rule = "rule1" },
    { name = "/tmp/abc_${time}.log", format = "04", rule = "rule2" },
]

# 文件更新检测周期（秒）
period = 30

# 服务本身的日志等级（生产环境请设置为info级别以上）
log_level = "debug"

alarm_url = "https://oapi.dingtalk.com/robot/send?access_token=1234567890"
```

1. 带日期的日志，使用`${time}`替换日期，并提供`format`。`abc_20210915.log`可以定义为`{name="abc_${time}.log", format="20060102"}`，`format`写法遵循golang时间format。

2. 钉钉报警器`&DingDingAlerter{url, NewTokenBucketLimiter(60, 8)}`默认限制60秒内可以发8条消息。如果拿不到令牌，则会将消息打到错误日志中。

   

