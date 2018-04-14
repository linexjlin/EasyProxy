package schedule

const (
	PollName   = "poll"
	IpHashName = "iphash"
	RandomName = "random"
	BackupName = "backup"
)

var registry = make(map[string]Strategy)

type Strategy interface {
	Init()
	Choose(client string, servers []string) string
}

func init() {
	registry[PollName] = new(Poll)
	registry[IpHashName] = new(IpHash)
	registry[RandomName] = new(Random)
	registry[BackupName] = new(Backup)
}

func GetStrategy(name string) Strategy {
	return registry[name]
}
