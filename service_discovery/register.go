package goDiscovery

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 参考https://blog.csdn.net/blogsun/article/details/102861648

const (
	EtcdDialTimeout   = 5 * time.Second
	EtcdConnTimeout   = 3 * time.Second
	KeepAliveInterval = time.Minute
	LeaseTTL          = 3
)

type KeyEncode func(prefix, key string) string

type ServiceRegister struct {
	clientv3.Client
	client          *clientv3.Client
	leaseID         clientv3.LeaseID
	keepAliveChan   <-chan *clientv3.LeaseKeepAliveResponse
	keepAliveCancel context.CancelFunc
	prefix          string
}

func NewServiceRegister(config clientv3.Config, ttl int64, prefix string) (*ServiceRegister, error) {
	client, err := clientv3.New(config)
	if err != nil {
		log.Println("NewServiceRegister", "init etcd client failed", err.Error())
		return nil, err
	}
	r := &ServiceRegister{client: client, prefix: prefix}
	// 设置租约并保活
	if err := r.SetLease(ttl); err != nil {
		return nil, err
	}
	go r.ListenLeaseResp()
	return r, nil
}

func (r *ServiceRegister) DefaultRegister() (err error) {
	k, v, err := r.DefaultKV()
	if err != nil {
		return
	}
	return r.RegisterKV(k, v)
}

func (r *ServiceRegister) concatPrefixKey(prefix, key string) string {
	return fmt.Sprintf("%s%s", prefix, key)
}

func (r *ServiceRegister) RegisterKV(k, v string) (err error) {
	err = r.PutKV(r.concatPrefixKey(r.prefix, k), v)
	if err != nil {
		log.Println("ServiceRegister.RegisterKV", "PutKV failed", err.Error(), k, v)
		return
	}
	log.Println("Register Service OK", k, v)
	return
}

func (r *ServiceRegister) SetLease(ttl int64) (err error) {
	// 设置租约
	ctx, cancel := context.WithTimeout(context.Background(), EtcdConnTimeout)
	defer cancel()
	leaseResp, err := r.client.Lease.Grant(ctx, ttl)
	if err != nil {
		log.Println("ServiceRegister.SetLease", "Lease Grant failed", err.Error())
		return
	}

	// 续约租期
	ctx, cancelFunc := context.WithCancel(context.Background())
	leaseRespChan, err := r.client.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		log.Println("ServiceRegister.SetLease", "KeepAlive return error", err.Error())
		return
	}

	r.leaseID = leaseResp.ID
	r.keepAliveCancel = cancelFunc
	r.keepAliveChan = leaseRespChan

	log.Println("Lease id:", r.leaseID)
	return
}

func (r *ServiceRegister) PutKV(k, v string) error {
	ctx, cancel := context.WithTimeout(context.Background(), EtcdConnTimeout)
	defer cancel()
	_, err := r.client.Put(ctx, k, v, clientv3.WithLease(r.leaseID))
	if err != nil {
		return err
	}
	return nil
}

func (r *ServiceRegister) ListenLeaseResp() {
	for resp := range r.keepAliveChan {
		if resp == nil {
			log.Println("cancel lease")
		}
	}
}

func (r *ServiceRegister) Close() error {
	// 取消续约
	r.keepAliveCancel()

	// 释放租约
	if _, err := r.client.Revoke(context.Background(), r.leaseID); err != nil {
		return err
	}
	log.Println("revoke lease ok", r.leaseID)
	return r.client.Close()
}

// key: hostname; value：timestamp
func (r *ServiceRegister) DefaultKV() (key, value string, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ServiceRegister.DefaultKV", "get hostname failed", err.Error())
		return
	}
	key = hostname
	value = strconv.FormatInt(time.Now().Unix(), 10)
	return
}

func DefaultServiceRegister(endpoint, user, password, prefix string) *ServiceRegister {
	conf := clientv3.Config{
		Endpoints:            []string{endpoint},
		DialTimeout:          EtcdDialTimeout,
		DialKeepAliveTime:    KeepAliveInterval,
		DialKeepAliveTimeout: EtcdDialTimeout,
		Username:             user,
		Password:             password,
	}
	// 实例化注册服务对象
	r, err := NewServiceRegister(conf, LeaseTTL, prefix)
	if err != nil {
		log.Fatal(err)
	}
	// 在指定prefix下注册服务
	if err = r.DefaultRegister(); err != nil {
		log.Fatal(err)
	}
	return r
}
