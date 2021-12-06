package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gobs/cmd"

	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/setting"
	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/util"
)

// cli所需的结构体
type pepperCmd struct {
	cmd.Cmd
	conn net.Conn
}

// cli实例command
var command = new(pepperCmd)
var address string
var port int

func init() {
	flag.StringVar(&address, "node", setting.DefaultAddress, "node address")
	flag.IntVar(&port, "port", setting.DefaultPort, "port")
	flag.Parse()
}

func main() {
	initCommandLine()
	registerCmd()
	command.CmdLoop()
}

// 初始化command实例
func initCommandLine() {
	command.Prompt = "pepper-cli> "
	command.PreLoop = func() {
		fmt.Println("Welcome pepperCache-client!")
		// 和server建立tcp连接
		c, e := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
		if e != nil {
			log.Panicln(e)
		}
		command.conn = c
	}
	command.Init()
}

func (c *pepperCmd) Send2Server(msg string, conn net.Conn) (n int, err error) {
	p, e := util.EncodeCmd(msg)
	fmt.Println("after", string(p))
	if e != nil {
		return 0, e
	}
	n, err = conn.Write(p)
	return n, err
}

func (c *pepperCmd) ReadServer(conn net.Conn) string {
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		serverAddrStr := conn.LocalAddr().String()
		return "can not connect to " + serverAddrStr
	}
	resp, er := util.DecodeFromBytes(buff)

	if n == 0 {
		return "nil"
	} else if er == nil {
		return string(resp.Value)
	} else {
		return "err server response"
	}
}
func checkError(err error) {
	if err != nil {
		log.Println("err ", err.Error())
		os.Exit(1)
	}
}

// 注册可用命令
func registerCmd() {
	command.Add(cmd.Command{
		Name: "set",
		Help: `set command`,
		Call: func(line string) bool { return addCommand("set", line) }})

	command.Add(cmd.Command{
		Name: "get",
		Help: `get command`,
		Call: func(line string) bool { return addCommand("get", line) }})

	command.Add(cmd.Command{
		Name: "del",
		Help: `del command`,
		Call: func(line string) bool { return addCommand("del", line) }})
}

func addCommand(name, line string) (stop bool) {
	msg := name + " " + line
	command.Send2Server(msg, command.conn)
	ret := command.ReadServer(command.conn)
	fmt.Println(ret)
	return
}
