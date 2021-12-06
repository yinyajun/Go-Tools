package core

import (
	"github.com/yinyajun/Golang-Toys/pepper_cache/cache/server/util"
)

const OK = "OK"
const NIL = "nil"

func addReplyStatus(c *Client, s string) {
	r := util.NewString([]byte(s))
	addReplyString(c, r)
}

func addReplyError(c *Client, s string) {
	r := util.NewError([]byte(s))
	addReplyString(c, r)
}

func addReplyString(c *Client, r *util.Resp) {
	if ret, err := util.EncodeToBytes(r); err == nil {
		c.Buf = ret
	}
}
