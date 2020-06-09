package scheduler

// ServerConfiguration contains configuration for
// one server and keys for the json parsing.
type ServerConfiguration struct {
	Name           string   `json:"name"`
	Address        string   `json:"address"`
	Port           uint     `json:"port"`
	Type           string   `json:"type"`
	SupportedTasks []string `json:"supportedTasks"`
}

// Config contains a list of server configurations
// for the scheduler with keys for json parsing.
type Config struct {
	AvailableServers []ServerConfiguration `json:"availableServers"`
}
