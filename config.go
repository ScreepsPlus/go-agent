package main

type agentConfig struct {
	Servers []agentConfigServer
}
type agentConfigServer struct {
	Server   string
	Shards   []string
	Segments []int
	Memory   string
	Console  struct {
		Prefix    string
		Seperator string
	}
	Prefix           string
	Interval         int
	ScreepsplusToken string
}
