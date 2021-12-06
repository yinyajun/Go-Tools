/*
* @Author: Yajun
* @Date:   2021/12/5 16:58
 */

package abtest

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type Storage interface {
	Init()                                       // 初始化配置
	Watch()                                      // 监控配置改动
	RegisterFunc(key string, f UpdateFunc) error // 注册更新函数
}

type UpdateFunc func(key string, data []byte) error

type LocalFile struct {
	dirname   string
	updateMap map[string]UpdateFunc
	watcher   *fsnotify.Watcher
}

func NewLocalFile(dirname string) *LocalFile {
	s := &LocalFile{dirname: dirname}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(dirname)
	if err != nil {
		log.Fatal(err)
	}
	s.watcher = watcher
	return s
}

func (l *LocalFile) RegisterFunc(key string, f UpdateFunc) error {
	if l.updateMap == nil {
		l.updateMap = make(map[string]UpdateFunc)
	}
	if _, ok := l.updateMap[key]; ok {
		return fmt.Errorf("duplicated updateMap key： %s", key)
	}
	if f == nil {
		return fmt.Errorf("updateFunc is nil")
	}
	l.updateMap[key] = f
	return nil
}

func (l *LocalFile) key(filename string) string {
	basename := path.Base(filename)
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}

func (l *LocalFile) Init() {
	files, err := ioutil.ReadDir(l.dirname)
	if err != nil {
		log.Panicln("LocalFile.Init", "ReadDir failed", err.Error())
	}
	var cnt int
	for _, f := range files {
		update, ok := l.updateMap[l.key(f.Name())]
		if !ok {
			continue
		}
		cnt++
		b, err := ioutil.ReadFile(path.Join(l.dirname, f.Name()))
		if err != nil {
			msg := fmt.Sprintf("ReadFile %s failed", f.Name())
			log.Panicln("LocalFile.Init", msg, err.Error())
		}
		key := l.key(f.Name())
		err = update(key, b)
		if err != nil {
			msg := fmt.Sprintf("update %s failed", key)
			log.Panicln("LocalFile.Init", msg, err.Error())
		}
		log.Printf("Init Domain %s ok\n", key)
	}
	if cnt != len(l.updateMap) {
		log.Panicln("LocalFile.Init", "Not all updateMap are executed", err.Error())
	}
}

func (l *LocalFile) Watch() {
	for {
		select {
		case e := <-l.watcher.Events:
			if e.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("Watch %s %s\n", e.Name, e.Op)
				b, err := ioutil.ReadFile(e.Name)
				key := l.key(e.Name)
				update, ok := l.updateMap[key]
				if !ok {
					continue
				}
				if err != nil {
					msg := fmt.Sprintf("ReadFile %s failed", e.Name)
					log.Println("LocalFile.Watch", msg, err)
					continue
				}
				err = update(key, b)
				if err != nil {
					msg := fmt.Sprintf("update %s failed", key)
					log.Println("LocalFile.Watch", msg, err)
				}
				log.Printf("Update Domain %s ok\n", key)
			}
		case err := <-l.watcher.Errors:
			log.Println("LocalFile.Watch", "watch errors", err)
		}
	}
}

//type Etcd struct {
//	client          *clientv3.Client
//	EtcdConnTimeout time.Duration
//	EtcdPrefix      string
//}
//
//func (f *Etcd) Init() (*Domain, error) {
//	f.client.KV = namespace.NewKV(f.client.KV, f.EtcdPrefix)
//	f.client.Watcher = namespace.NewWatcher(f.client.Watcher, f.EtcdPrefix)
//
//	ctx, cancel := context.WithTimeout(context.Background(), f.EtcdConnTimeout)
//	resp, err := f.client.Get(ctx, "\x00", clientv3.WithFromKey())
//	defer cancel()
//	if err != nil {
//		log.Panicln("settings.initEtcdConf", "get etcd conf failed", err.Error())
//	}
//
//	data := make([]NamedData, 0)
//	for _, kv := range resp.Kvs {
//		data = append(data, NamedData{
//			Name: string(kv.Key),
//			Data: kv.Value,
//		})
//	}
//	return ParseNamedData(nil, data)
//}
