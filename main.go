package main

import "fmt"

func main() {
	cfg := parseConfig("sunny_5_skiers/config.json")
	fmt.Printf("%+v\n", cfg)
}
