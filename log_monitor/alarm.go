package log_monitor

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Alerter interface {
	Alarm(msg Alarm)
}

type TokenBucketLimiter struct {
	t      *time.Ticker
	bucket chan struct{}
	doneC  chan struct{}
	limit  int
	stop   sync.Once
}

// rate秒内能生成limit个令牌
func NewTokenBucketLimiter(rate, limit int) *TokenBucketLimiter {
	t := time.NewTicker(time.Duration(rate) * time.Second)
	doneC := make(chan struct{})
	bucket := make(chan struct{}, limit)
	tbl := &TokenBucketLimiter{
		t:      t,
		bucket: bucket,
		doneC:  doneC,
		limit:  limit,
		stop:   sync.Once{},
	}
	// 定时放置令牌
	tbl.asyncPutTokens()
	return tbl
}

func (limiter *TokenBucketLimiter) Allow() bool {
	select {
	case limiter.bucket <- struct{}{}:
		return true
	case <-limiter.doneC:
		return false
	default:
		return false
	}
}

func (limiter *TokenBucketLimiter) Close() {
	limiter.stop.Do(func() {
		close(limiter.doneC)
		close(limiter.bucket)
		limiter.t.Stop()
	})
}

// 一旦rate时间到了，就开始放置令牌
func (limiter *TokenBucketLimiter) asyncPutTokens() {
	go func() {
		for {
			select {
			case <-limiter.t.C:
				limiter.drain()
			case <-limiter.doneC:
				return
			}
		}
	}()
}

// 放置limit个令牌
func (limiter *TokenBucketLimiter) drain() {
	for i := 0; i < limiter.limit; i++ {
		select {
		case <-limiter.bucket:
		default:
			return
		}
	}
}

type Alarm struct {
	FileName string
	LineText string
	Rule     string
	Level    int
}

type DebugAlerter struct {
	limiter *TokenBucketLimiter
}

func NewDebugAlerter() *DebugAlerter {
	return &DebugAlerter{limiter: NewTokenBucketLimiter(10, 2)}
}

func (a *DebugAlerter) Alarm(msg Alarm) {
	if a.limiter.Allow() {
		DebugLogger.Println("[alarm]", msg.LineText, msg.Level, msg.FileName)
	} else {
		DebugLogger.Println("[unsent alarm]", msg.LineText, msg.Level, msg.FileName)
	}

}

// 钉钉报警器
type DingDingAlerter struct {
	webhookURL string
	limiter    *TokenBucketLimiter
}

func NewDingDingAlerter(url string) *DingDingAlerter {
	return &DingDingAlerter{url, NewTokenBucketLimiter(60, 8)}
}

func (a *DingDingAlerter) Alarm(msg Alarm) {
	hostname, _ := os.Hostname()
	title := "[%s][%s]\n"
	text := fmt.Sprintf("[file]: %s\n[rule]: %s\n[content]: %s\n", msg.FileName, msg.Rule, msg.LineText)
	switch msg.Level {
	case Error:
		title = fmt.Sprintf(title, hostname, "error")
	case Warning:
		title = fmt.Sprintf(title, hostname, "warning")
	}
	if a.limiter.Allow() {
		content := `{"msgtype": "text", "text": {"content": "` + title + text + `"}}`
		req, err := http.NewRequest("POST", a.webhookURL, strings.NewReader(content))
		if err != nil {
			ErrorLogger.Println("Alarm failed", msg)
		}
		client := &http.Client{}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		resp, err := client.Do(req)
		defer resp.Body.Close()
	} else { // 限流器限制不允许钉钉报警，将报警用error logger记录
		ErrorLogger.Println("[unsent alarm]", msg.LineText, msg.Level, msg.FileName)
	}
}
