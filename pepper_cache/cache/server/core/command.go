package core

import (
	"time"
)

func lookupKey(db *PepperDb, key *PepperObject) (ret *PepperObject) {
	if o, ok := db.Dict.Get(key.Ptr.(string)); ok {
		return o.(*PepperObject)
	}
	return nil
}

// Get 命令
func GetCommand(c *Client, s *Server) {
	if c.Argc != 2 {
		addReplyError(c, "(error) ERR wrong number of arguments for 'get' command")
		return
	}
	objKey := c.Argv[1]
	v := lookupKey(c.Db, objKey)
	if v != nil {
		addReplyStatus(c, v.Ptr.(string))
	} else {
		addReplyStatus(c, NIL)
	}
}

// Set 命令
func SetCommand(c *Client, s *Server) {
	if c.Argc != 3 {
		addReplyError(c, "(error) ERR wrong number of arguments for 'set' command")
		return
	}
	objKey := c.Argv[1]
	objValue := c.Argv[2]

	stringKey, ok1 := objKey.Ptr.(string)
	stringValue, ok2 := objValue.Ptr.(string)

	if !ok1 || !ok2 {
		addReplyError(c, NIL)
		return
	}

	c.Db.del(objKey, objValue)
	c.Db.Dict.Set(stringKey, CreateObject(OBJ_STRING, stringValue), 0)
	c.Db.add(objKey, objValue)
	s.Dirty++
	addReplyStatus(c, OK)
}

// Setex 命令
func SetexCommand(c *Client, s *Server) {
	if c.Argc != 4 {
		addReplyError(c, "(error) ERR wrong number of arguments for 'setex' command")
		return
	}
	objKey := c.Argv[1]
	objValue := c.Argv[2]
	objExpire := c.Argv[3]

	stringKey, ok1 := objKey.Ptr.(string)
	stringValue, ok2 := objValue.Ptr.(string)
	intExpire, ok3 := objExpire.Ptr.(time.Duration)

	if !ok1 || !ok2 || !ok3 {
		addReplyError(c, NIL)
		return
	}

	c.Db.del(objKey, objValue)
	c.Db.Dict.Set(stringKey, *CreateObject(OBJ_STRING, stringValue), intExpire*time.Second)
	c.Db.add(objKey, objValue)
	s.Dirty++
	addReplyStatus(c, OK)
}

// Del 命令
func DelCommand(c *Client, s *Server) {
	if c.Argc != 2 {
		addReplyError(c, "(error) ERR wrong number of arguments for 'del' command")
		return
	}
	objKey := c.Argv[1]
	v := lookupKey(c.Db, objKey)
	if v != nil {
		c.Db.del(objKey, v)
		c.Db.Dict.Del(objKey.Ptr.(string))
		s.Dirty++
		addReplyStatus(c, OK)
	} else {
		addReplyStatus(c, "nil")
	}
}

// todo:Expire 命令

// todo:TTL 命令

// todo:Exists 命令

// todo:Rename 命令

// todo:Setnx 命令

// 已注册的命令
func RegisteredCommands() map[string]*PepperCommand {
	getCommand := &PepperCommand{Name: "get", Proc: GetCommand}
	setCommand := &PepperCommand{Name: "set", Proc: SetCommand}
	setexCommand := &PepperCommand{Name: "setex", Proc: SetexCommand}
	delCommand := &PepperCommand{Name: "set", Proc: DelCommand}

	commands := map[string]*PepperCommand{
		"get":   getCommand,
		"set":   setCommand,
		"setex": setexCommand,
		"del":   delCommand,
	}
	return commands
}
