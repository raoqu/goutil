package shell

type AliveCheckConfig struct {
	Ping string
}

func (c *Command) UpdateStatus() CommandStatus {
	if c.Attached {
		if len(c.AliveConfig.Ping) > 0 {
			psPing(c)
		}
	}

	return c.Status
}

func psPing(c *Command) bool {
	stat := ShellStat{}
	contains := stat.CheckOutput("ps -ef", c.AliveConfig.Ping)
	c.extBuffer = stat.psoutput
	if contains {
		c.Status = START
	} else {
		c.Status = UNKNOWN
	}
	return contains
}

func MapCommandStatus(status CommandStatus) string {
	switch status {
	case INIT:
		return "init"
	case FAIL:
		return "fail"
	case START:
		return "start"
	case COMPLETE:
		return "complete"
	default:
		return "unknown"
	}
}

func IsCommandStatusAlive(status CommandStatus) bool {
	switch status {
	case INIT:
	case START:
		return true
	default:
		return false
	}

	return false
}
