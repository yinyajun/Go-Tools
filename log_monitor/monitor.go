package log_monitor

import (
	"github.com/nxadm/tail"
	"io"
	"strings"
	"time"
)

const (
	timeFlag = "${time}"
)

var (
	seps = Seps{'\n'}
)

type Monitor struct {
	cfg        *Config
	alerts     []Alerter
	manager    []Cell
	failedChan chan Cell
	alarmChan  chan Alarm
	rules      map[string]*Trie
}

type Cell struct {
	int             // 在manager中的index
	file File       // 配置中的文件信息
	tail *tail.Tail // 实际打开的文件
}

func NewMonitor(cfg *Config, alerts []Alerter) *Monitor {
	m := &Monitor{
		cfg:        cfg,
		alarmChan:  make(chan Alarm, 100),
		failedChan: make(chan Cell, 20),
		manager:    make([]Cell, len(cfg.Files)),
		alerts:     alerts,
		rules:      make(map[string]*Trie),
	}
	// init rules
	for _, r := range m.cfg.Rules {
		m.rules[r.Name] = InitTree(r.ExtractDict())
	}
	return m
}

func (m *Monitor) MonitFile(c Cell, time time.Time) {
	fileName := GetFileName(c.file, time)
	seekInfo := tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	t, err := tail.TailFile(fileName, tail.Config{Follow: true, Poll: true, ReOpen: true, Location: &seekInfo, Logger: DebugLogger})
	if err != nil {
		WarningLogger.Println(c.file, err)
		m.failedChan <- c
		DebugLogger.Println("[monitFile] put into failed chan", c.file.Name)
		return
	}
	// file exists
	c.tail = t
	go m.ExecuteTail(c, time)
}

func (m *Monitor) ExecuteTail(c Cell, time time.Time) {
	defer c.tail.Stop()
	defer c.tail.Cleanup()
	m.manager[c.int] = c
	DebugLogger.Println("[ExecuteTail] put", c.tail.Filename)

	rule, exist := m.rules[c.file.Rule]
	if !exist {
		ErrorLogger.Printf("rule %s not exists", c.file.Rule)
		return
	}

	for line := range c.tail.Lines {
		DebugLogger.Println(">>", line.Text)
		res := rule.Match([]rune(line.Text), seps)
		if len(res[Ignored]) > 0 {
			continue
		}
		if len(res[Error]) > 0 {
			m.alarmChan <- Alarm{c.tail.Filename, line.Text, c.file.Rule, Error}
			continue
		}
		if len(res[Warning]) > 0 {
			m.alarmChan <- Alarm{c.tail.Filename, line.Text, c.file.Rule, Warning}
			continue
		}
	}
	// to avoid race map manager
	DebugLogger.Println("[ExecuteTail] tail lines ends, ", c.tail.Filename, c.tail.Err())
	m.failedChan <- c
}

func (m *Monitor) ExecuteFail(time time.Time) {
	for {
		select {
		case c := <-m.failedChan:
			DebugLogger.Println("[ExecuteFail] meets", c)
			go m.MonitFile(c, time)
		default:
			DebugLogger.Println("[ExecuteFail] ends")
			return
		}
	}
}

func (m *Monitor) Update(triggerTime time.Time) {
	DebugLogger.Println("--------------------------------------------------------------")
	DebugLogger.Println("[Update] begins")

	for _, c := range m.manager {
		if c.tail == nil {
			continue
		}
		DebugLogger.Println("[Update] meets", c.tail.Filename)
		if NeedUpdate(c.file, c.tail, triggerTime) {
			DebugLogger.Println(">> need to update ", c.tail.Filename)
			c.tail.Stop()
			c.tail.Cleanup()
		}
	}
	time.Sleep(1 * time.Second)
	go m.ExecuteFail(triggerTime)
}

func (m *Monitor) WaitForAlarm() {
	for {
		select {
		case msg := <-m.alarmChan:
			m.Alarm(msg)
		}
	}
}

func (m *Monitor) Alarm(msg Alarm) {
	for _, a := range m.alerts {
		a.Alarm(msg)
	}
}

func (m *Monitor) Monit() {
	now := time.Now()
	for idx, f := range m.cfg.Files {
		go m.MonitFile(Cell{idx, f, nil}, now)
	}
	go m.WaitForAlarm()
	timer := time.NewTicker(time.Duration(m.cfg.Period) * time.Second)
	for {
		triggerTime := <-timer.C // block util trigger
		m.Update(triggerTime)
	}
}

func NeedUpdate(f File, t *tail.Tail, cur time.Time) bool {
	name := GetFileName(f, cur)
	if name == t.Filename {
		return false
	}
	return true
}

func GetFileName(f File, t time.Time) string {
	if f.Format == "" {
		return f.Name
	}
	return strings.Replace(f.Name, timeFlag, t.Format(f.Format), -1)
}
