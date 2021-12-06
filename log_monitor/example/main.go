package main

import (
	"log"
	monitor "log_monitor"
)

var (
	m *monitor.Monitor
)

func init() {
	// 读取配置
	cfg, err := monitor.NewConfigWithFile("example/log.toml")
	if err != nil {
		log.Panicln(err)
	}
	// 初始化日志
	monitor.InitLogger(cfg.LogLevel)
	// 初始化报警器
	alerts := []monitor.Alerter{
		monitor.NewDebugAlerter(),
		monitor.NewDingDingAlerter(cfg.AlarmURL),
	}
	// 初始化日志监视器
	m = monitor.NewMonitor(cfg, alerts)
}

func main() {
	m.Monit()
}
