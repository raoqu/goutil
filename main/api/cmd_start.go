package api

type StartParam struct {
	Dir     string   `json:"dir"`
	Command string   `json:"command"`
	Params  []string `json:"params"`
}

type StartResult struct {
	ProcessID int64  `json:"process_id"`
	Error     string `json:"error"`
	Success   bool   `json:"success"`
}

func StartCommand(param StartParam) (StartResult, error) {
	return StartResult{
		ProcessID: 0,
		Success:   true,
	}, nil
}
