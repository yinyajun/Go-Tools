package setting

const (
	DefaultCacheSlotNum = 256
	DefaultAofFile      = "./pepper.aof"
	DefaultPort         = 9004
	DefaultAdminPort    = 18088
	DefaultAddress      = "127.0.0.1"
)

var (
	Port      int
	AdminPort int
	Address   string
	Cluster   string
)
