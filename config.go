package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Servers map[string]struct {
		Host string
		Port string
		Secure bool		
		Token string
		Username string
		Password string
  	Shards []string
		Segments []int
	  Memory []string
	  Console struct {
	    Prefix string
	    Seperator string
	  }
	}
	Screepsplus struct {
		Url string
		Token string
	}
}

func ReadConfig(file string) *Config {
	var config Config
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Config failed to open #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Config Unmarshal %v", err)
	}
	return &config
}