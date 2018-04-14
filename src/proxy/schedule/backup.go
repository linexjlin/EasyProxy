package schedule

type Backup struct {
}

func (strategy *Backup) Init() {}

func (strategy *Backup) Choose(client string, servers []string) string {
	url := servers[0]
	return url
}
