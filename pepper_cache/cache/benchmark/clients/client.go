package clients

type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}

type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}

func New(typ, server, port string) Client {
	if typ == "redis" {
		return newRedisClient(server, port)
	}
	if typ == "tcp" {
		return newTCPClient(server, port)
	}
	panic("unknown clients type " + typ)
}
