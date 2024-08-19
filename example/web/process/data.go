package process

type Command struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Config struct {
	Uuid    string `json:"uuid"`
	Dir     string `json:"dir"`
	Command string `json:"command"`
	Ping    string `json:"ping"`
}

type Stat struct {
	Status string   `json:"status"`
	Output []string `json:"output"`
}
