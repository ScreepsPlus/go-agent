package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/screepers/go-screeps/config"
	"github.com/screepers/go-screeps/screeps"
)

func main() {
	aconf := &agentConfig{
		Servers: []agentConfigServer{
			{
				Server:           "main",
				Segments:         []int{1},
				Memory:           "stats",
				Interval:         10,
				ScreepsplusToken: "S+ Token",
				Shards: []string{
					"shard1",
					"shard2",
				},
			},
			{
				Server:           "screepsplus",
				Segments:         []int{1, 2},
				Memory:           "stats",
				Interval:         10,
				ScreepsplusToken: "S+ Token",
				Shards: []string{
					"screepsplus1",
				},
			},
		},
	}
	conf := config.NewConfig()
	for _, server := range aconf.Servers {
		go runServer(conf.Servers[server.Server], server)
	}
	select {}
}

func runServer(conf config.ServerConfig, server agentConfigServer) {
	client := screeps.NewClient(conf)
	//client.SetDebug(true)
	for {
		start := time.Now()
		stats := make([]Stat, 0)
		sources := make([]string, 0)
		if server.Memory != "" {
			sources = append(sources, "memory")
		}
		if len := len(server.Segments); len > 0 {
			sources = append(sources, "segments")
		}
		for _, shard := range server.Shards {
			if server.Memory != "" {
				mem, err := client.GetMemory("stats", shard)
				if err != nil {
					log.Printf("%v", err)
					return
				}
				nstats := processStats(mem, server.Prefix)
				stats = append(stats, nstats...)
			}
			if len(server.Segments) == 1 {
				segment, err := client.GetMemorySegment(server.Segments[0], shard)
				if err != nil {
					log.Printf("%v", err)
					return
				}
				nstats := processStats(segment, server.Prefix)
				stats = append(stats, nstats...)
			}
			if len(server.Segments) > 1 {
				segments, err := client.GetMemorySegments(server.Segments, shard)
				if err != nil {
					log.Printf("%v", err)
					return
				}
				for _, segment := range segments {
					nstats := processStats(&segment, server.Prefix)
					stats = append(stats, nstats...)
				}
			}
		}
		elapsed := time.Since(start)
		log.Printf("[%s] Fetched %d stats in %dms from sources: [%s] on shards: [%s]", server.Server, len(stats), elapsed/time.Millisecond, strings.Join(sources, ","), strings.Join(server.Shards, ","))
		<-time.After((time.Duration(server.Interval) * time.Second) - elapsed)
	}
}

func processStats(mem *screeps.GetMemoryResponse, prefix string) []Stat {
	var rawStats interface{}
	rawStats = make(map[string]interface{})
	mem.Parse(&rawStats)
	stats := flattenJSON(rawStats, prefix)
	return stats
}

// Stat - a single stat
type Stat struct {
	Key    string
	Value  float64
	Labels map[string]string
}

func flattenJSON(data interface{}, key string) []Stat {
	ret := make([]Stat, 0)
	switch v := data.(type) {
	case nil: // Ignore nil
	case string: // Ignore strings
	case float64:
		stat := Stat{
			Key:   key,
			Value: v,
		}
		ret = append(ret, stat)
	case []interface{}:
		for i, vv := range v {
			subKey := strconv.Itoa(i)
			if key != "" {
				subKey = fmt.Sprintf("%s.%s", key, subKey)
			}
			res := flattenJSON(vv, subKey)
			ret = append(ret, res...)
		}
	case map[string]interface{}:
		for i, vv := range v {
			subKey := i
			if key != "" {
				subKey = fmt.Sprintf("%s.%s", key, subKey)
			}
			res := flattenJSON(vv, subKey)
			ret = append(ret, res...)
		}
	}
	return ret
}
