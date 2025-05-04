package main

import (
	"log"
)

func main() {
	cfg, err := parseConfig("sunny_5_skiers/config.json")
	if err != nil {
		log.Panic("Cannot load config")
	}

	_ = cfg
}
