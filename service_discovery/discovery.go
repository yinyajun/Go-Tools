package goDiscovery

import (
	"context"
	"fmt"
	"log"
	"os"

	c "discovery/consistent"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 参考https://blog.csdn.net/blogsun/article/details/102861648

type ServiceDiscovery struct {
	*c.Circle
	client          *clientv3.Client
	watchCancelFunc context.CancelFunc
	prefix          string
}

func NewServiceDiscovery(config clientv3.Config, prefix string, replicaNum int, hash string) (*ServiceDiscovery, error) {
	circle := c.NewCircle()
	circle.Hash = hash
	circle.ReplicaNum = replicaNum // 虚拟节点数目

	client, err := clientv3.New(config)
	if err != nil {
		log.Println("NewServiceDiscovery", "init etcd client failed", err.Error())
		return nil, err
	}
	d := &ServiceDiscovery{
		Circle: circle,
		client: client,
		prefix: prefix,
	}
	return d, nil
}

func (d *ServiceDiscovery) separateKey(key, prefix string) (string, error) {
	if len(key) <= len(prefix) {
		return "", fmt.Errorf("invalid key length")
	}
	return key[len(prefix):], nil
}

func (d *ServiceDiscovery) InitService() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), EtcdConnTimeout)
	defer cancel()

	resp, err := d.client.Get(ctx, d.prefix, clientv3.WithPrefix())
	if err != nil {
		log.Panicln("ServiceDiscovery.GetService", "get failed", err.Error())
		return
	}
	// 将结果存到hash环中
	members := []string{}
	for _, kv := range resp.Kvs {
		member, err := d.separateKey(string(kv.Key), d.prefix)
		if err != nil {
			continue
		}
		members = append(members, member)
		d.Add(member)
	}
	log.Println("Current Members:", d.GetService())
	return
}

func (d *ServiceDiscovery) WatchService() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	d.watchCancelFunc = cancelFunc
	for resp := range d.client.Watch(ctx, d.prefix, clientv3.WithPrefix()) {
		for _, kv := range resp.Events {
			member, err := d.separateKey(string(kv.Kv.Key), d.prefix)
			if err != nil {
				continue
			}
			switch kv.Type {
			case mvccpb.PUT:
				log.Println("Add Member:", member)
				d.Add(member)
			case mvccpb.DELETE:
				log.Println("Delete Member:", member)
				d.Del(member)
			}
		}
	}
}

func (d *ServiceDiscovery) GetService() []string { return d.Members() }

// match hostname
func (d *ServiceDiscovery) MatchHost(key string) bool {
	hostname, err1 := os.Hostname()
	name, err2 := d.Allocate(key)
	if err1 != nil || err2 != nil {
		return false
	}
	return name == hostname
}

func (d *ServiceDiscovery) Close() error {
	d.watchCancelFunc()
	return d.client.Close()
}

func DefaultServiceDiscovery(endpoint, user, password, prefix string, replicaNum int, hash string) *ServiceDiscovery {
	conf := clientv3.Config{
		Endpoints:            []string{endpoint},
		DialTimeout:          EtcdDialTimeout,
		DialKeepAliveTime:    KeepAliveInterval,
		DialKeepAliveTimeout: EtcdDialTimeout,
		Username:             user,
		Password:             password,
	}
	d, err := NewServiceDiscovery(conf, prefix, replicaNum, hash)
	if err != nil {
		log.Fatal(err)
	}

	if err = d.InitService(); err != nil {
		log.Fatal(err)
	}
	go d.WatchService()
	return d
}
