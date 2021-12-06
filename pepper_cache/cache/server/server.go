package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/admin"
	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/core"
	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/dict"
	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/setting"
)

// 服务端实例
var pepper = new(core.Server)

func main() {
	//读取配置
	setting.Init()
	//初始化服务端实例
	initServer()
	// 服务端通过tcp监听，负责kv的操作
	go pepper.TcpListen(handle)
	// 服务端通过http监听，负责集群和节点状态信息查询
	ad := admin.NewAdminServer(pepper)
	go ad.HttpListen()
	select {}
}

// 初始化服务端实例
func initServer() {
	//实例化cluster.Node
	n, e := core.NewNode(setting.Address, setting.Cluster)
	if e != nil {
		panic(e)
	}
	//初始化Server实例
	pepper.Node = n
	pepper.Pid = os.Getpid()
	pepper.DbNum = 16
	pepper.Port = int32(setting.Port)
	initDb()
	pepper.Start = time.Now().UnixNano() / 1000000
	pepper.AofFileName = setting.DefaultAofFile
	pepper.Commands = core.RegisteredCommands()
	//加载持久化文件
	LoadData()
}

// 初始化db
func initDb() {
	pepper.Db = make([]*core.PepperDb, pepper.DbNum)
	for i := 0; i < pepper.DbNum; i++ {
		pepper.Db[i] = new(core.PepperDb)
		pepper.Db[i].Dict = dict.NewCache(time.Minute, setting.DefaultCacheSlotNum)
	}
}

// 加载aof持久化文件
func LoadData() {
	c := pepper.CreateClient()
	c.FakeFlag = true
	pros := core.ReadAof(pepper.AofFileName)
	for _, v := range pros {
		c.QueryBuf = string(v)
		err := c.ProcessInputBuffer()
		if err != nil {
			log.Println("ProcessInputBuffer err", err)
		}
		pepper.SyncProcessCommand(c)
	}
}

// 处理请求: 异步操作（确保已经实例化pepper）
func handle(conn net.Conn) {
	c := pepper.CreateClient()
	// 创建二维channel来存放result
	resultCh := make(chan chan core.Result, 5000)
	c.ResultCh = resultCh
	defer close(resultCh)

	go responseConn(conn, c)

	for {
		//同步读取conn的内容
		err := c.ReadQueryFromClient(conn)
		if err != nil {
			//log.Println("readQueryFromClient err", err)
			return
		}
		err = c.ProcessInputBuffer()
		if err != nil {
			//log.Println("ProcessInputBuffer err", err)
			return
		}
		//执行操作
		pepper.ProcessCommand(c)
	}
}

// 异步返回响应
func responseConn(conn net.Conn, c *core.Client) {
	defer conn.Close()
	for ch := range c.ResultCh {
		r := <-ch
		conn.Write(r)
	}
}
