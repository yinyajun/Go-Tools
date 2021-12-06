package log_monitor

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

const (
	Ignored = 1
	Warning = 2
	Error   = 3
)

type Rule struct {
	Name         string   `toml:"name"`
	IgnoredWords []string `toml:"ignored_words"`
	ErrorWords   []string `toml:"error_words"`
	WarningWords []string `toml:"warning_words"`
}

type File struct {
	Name   string `toml:"name"`
	Format string `toml:"format"`
	Rule   string `toml:"rule"`
}

type Config struct {
	Rules    []Rule `toml:"rules"`
	Files    []File `toml:"files"`
	Period   int64  `toml:"period"`
	LogLevel string `toml:"log_level"`
	AlarmURL string `toml:"alarm_url"`
}

func NewConfigWithFile(name string) (*Config, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewConfig(string(data))
}

func NewConfig(data string) (*Config, error) {
	var c Config
	_, err := toml.Decode(data, &c)
	if err != nil {
		return nil, err
	}
	if err := c.Check(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) Check() error {
	rules := map[string]struct{}{}
	var dict map[int][]string
	for _, r := range c.Rules {
		if r.Name == "" {
			return fmt.Errorf("empty rule name in %s", r)
		}
		dict = r.ExtractDict()
		if len(dict[Error]) == 0 && len(dict[Warning]) == 0 {
			return fmt.Errorf("both error_words and warning_words are empty in %s", r.Name)
		}
		rules[r.Name] = struct{}{}
	}

	for _, f := range c.Files {
		if f.Name == "" {
			return fmt.Errorf("empty file name in %s", f)
		}
		if _, exist := rules[f.Rule]; !exist {
			return fmt.Errorf("rule unavailable in %s", f.Name)
		}
	}
	return nil
}

func (r *Rule) ExtractDict() map[int][]string {
	return map[int][]string{
		Ignored: r.IgnoredWords,
		Error:   r.ErrorWords,
		Warning: r.WarningWords,
	}
}
