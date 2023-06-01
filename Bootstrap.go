package main

import (
	"fmt"
	"log"
	"os"
	"pluginengine/utils"
	"strings"
)

// main takes:
// 1) IP of network device
// 2) function name (discovery/collect)
// 3) Metric Group(s)(Only if 2nd parameter is collect)
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Missing Command Line Arguments")
			fmt.Println("Error: ", err)
		}
	}()

	var ip = os.Args[1]
	var action = os.Args[2]
	if strings.EqualFold(action, "discovery") {
		status, err := utils.Discovery(ip)
		fmt.Println(map[string]any{"status": status, "error": err})
	} else if strings.EqualFold(action, "collect") {

	} else {
		panic("Unknown Function")
	}
}
