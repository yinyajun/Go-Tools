package core

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/dict"
	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/util"
)

//Client 与服务端连接之后即创建一个Client结构
type Client struct {
	Cmd      *PepperCommand
	Argv     []*PepperObject
	Argc     int
	Db       *PepperDb
	Buf      []byte
	QueryBuf string
	ResultCh chan chan Result
	FakeFlag bool
}

// Server 服务端实例结构体
type Server struct {
	Node             Node
	Db               []*PepperDb
	DbNum            int
	Start            int64
	Port             int32
	AofFileName      string
	RdbFileNmae      string
	SystemMemorySize int32
	Clients          int32
	Pid              int
	Commands         map[string]*PepperCommand
	Dirty            int64
	AofBuf           []string
}

// db 结构体
type PepperDb struct {
	Stat
	Dict dict.Dict
	ID   int32
}

//Command 命令结构
type PepperCommand struct {
	Name string
	Proc cmdFunc
}

//命令函数指针
type cmdFunc func(c *Client, s *Server)

//Result结构体
type Result []byte

//conn处理函数指针
type Handler func(net.Conn)

// Server通过tcp监听端口
func (s *Server) TcpListen(handler Handler) {
	a := fmt.Sprintf("%s:%d", s.Node.Addr(), s.Port)
	l, err := net.Listen("tcp", a)

	if err != nil {
		log.Print("listen err ")
	}
	log.Println("Cache Server runs...")
	log.Println("Cache Port:", s.Port)
	log.Println("Cache PID:", s.Pid)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		log.Println("local:", conn.LocalAddr(), "remote:", conn.RemoteAddr())
		go handler(conn)
	}
}

func (s *Server) SyncProcessCommand(c *Client) {
	CmdName, ok := c.Argv[0].Ptr.(string)
	if !ok {
		log.Println("error cmd")
		os.Exit(1)
	}
	cmd := lookupCommand(CmdName, s)
	if cmd == nil {
		addReplyError(c, fmt.Sprintf("(error) ERR unknown command '%s'", CmdName))
		return
	}
	c.Cmd = cmd
	call(c, s)
}

// 执行命令
func (s *Server) ProcessCommand(c *Client) {
	// 同步检查命令
	CmdName, ok := c.Argv[0].Ptr.(string)
	if !ok {
		log.Println("error cmd")
		os.Exit(1)
	}
	cmd := lookupCommand(CmdName, s)
	if cmd == nil {
		addReplyError(c, fmt.Sprintf("(error) ERR unknown command '%s'", CmdName))
		return
	}
	c.Cmd = cmd
	// 异步操作
	ch := make(chan Result)
	c.ResultCh <- ch
	// todo 并发的协程不应该依赖client的成员，因为无法保证client的成员不变(目前用锁解决)
	go func() {
		// 判断命令是否需要在该server上执行
		if ok = s.shouldProcess(c); ok {
			call(c, s)
		}
		ch <- Result(c.Buf)
	}()
}

// 是否需在该server上处理（不需要处理时，给出错误result）
func (s *Server) shouldProcess(c *Client) bool {
	// 如果是加载aof的client
	if c.FakeFlag {
		return true
	}
	//不包含key的命令
	if c.Argc < 2 {
		return true
	}
	keyName, ok := c.Argv[1].Ptr.(string)
	if !ok {
		addReplyError(c, fmt.Sprintf("(error) ERR invalid key '%s'", keyName))
		return false
	}
	addr, ok := s.Node.ShouldProcess(keyName)
	// 不需要在该server上处理该key
	if !ok {
		addReplyError(c, fmt.Sprintf("(error) redirect %s", addr))
		return false
	}
	return true
}

// 从注册命令中查找命令
func lookupCommand(name string, s *Server) *PepperCommand {
	if cmd, ok := s.Commands[name]; ok {
		return cmd
	}
	return nil
}

// 真正调用命令
func call(c *Client, s *Server) {
	dirty := s.Dirty
	c.Cmd.Proc(c, s)
	dirty = s.Dirty - dirty
	if dirty > 0 && !c.FakeFlag {
		AppendToFile(s.AofFileName, c.QueryBuf)
	}
}

// 创建client记录当前连接
func (s *Server) CreateClient() (c *Client) {
	c = new(Client)
	c.Db = s.Db[0]
	c.QueryBuf = ""
	return c
}

// 从当前连接中读取请求
func (c *Client) ReadQueryFromClient(conn net.Conn) (err error) {
	buff := make([]byte, 512)
	n, err := conn.Read(buff)

	if err != nil {
		log.Println("conn.Read err!=nil", err, n, conn.RemoteAddr())
		conn.Close()
		return err
	}
	c.QueryBuf = string(buff)
	return nil
}

// 处理请求信息
func (c *Client) ProcessInputBuffer() error {
	decoder := util.NewDecoder(bytes.NewReader([]byte(c.QueryBuf)))
	if resp, err := decoder.DecodeMultiBulk(); err == nil {
		c.Argc = len(resp)
		c.Argv = make([]*PepperObject, c.Argc)
		for k, s := range resp {
			c.Argv[k] = CreateObject(OBJ_STRING, string(s.Value))
		}
		return nil
	}
	return errors.New("ProcessInputBuffer failed")
}
