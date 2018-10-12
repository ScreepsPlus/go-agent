package main

import "log"

func main() {
	conf := ReadConfig("config.yaml")
	log.Printf("Config: %v", conf)
}