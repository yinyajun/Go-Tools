package clients

import (
	"errors"
	"fmt"
	"net"

	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/util"
)

type tcpClient struct {
	net.Conn
}

func newTCPClient(server, port string) *tcpClient {
	c, e := net.Dial("tcp", server+":"+port)
	if e != nil {
		panic(e)
	}
	return &tcpClient{c}
}

func (c *tcpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		msg := fmt.Sprintf("%s %s", cmd.Name, cmd.Key)
		c.Send2Server(msg)
		cmd.Value, cmd.Error = c.ReadServer()
		return
	}
	if cmd.Name == "set" {
		msg := fmt.Sprintf("%s %s %s", cmd.Name, cmd.Key, cmd.Value)
		c.Send2Server(msg)
		_, cmd.Error = c.ReadServer()
		return
	}
	if cmd.Name == "del" {
		msg := fmt.Sprintf("%s %s", cmd.Name, cmd.Key)
		c.Send2Server(msg)
		_, cmd.Error = c.ReadServer()
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

func (c *tcpClient) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		if cmd.Name == "get" {
			msg := fmt.Sprintf("%s %s", cmd.Name, cmd.Key)
			c.Send2Server(msg)
		}
		if cmd.Name == "set" {
			msg := fmt.Sprintf("%s %s %s", cmd.Name, cmd.Key, cmd.Value)
			c.Send2Server(msg)
		}
		if cmd.Name == "del" {
			msg := fmt.Sprintf("%s %s", cmd.Name, cmd.Key)
			c.Send2Server(msg)
		}
	}
	for _, cmd := range cmds {
		cmd.Value, cmd.Error = c.ReadServer()
	}
}

func (c *tcpClient) Send2Server(msg string) (n int, err error) {
	p, e := util.EncodeCmd(msg)
	if e != nil {
		return 0, e
	}
	n, err = c.Write(p)
	return n, err
}

func (c *tcpClient) ReadServer() (string, error) {
	buff := make([]byte, 1024)
	n, err := c.Read(buff)
	if err != nil {
		return "", err
	}
	resp, er := util.DecodeFromBytes(buff)
	if n == 0 {
		return "", nil
	} else if er == nil {
		return string(resp.Value), nil
	} else {
		return "", errors.New("err server response")
	}
}
