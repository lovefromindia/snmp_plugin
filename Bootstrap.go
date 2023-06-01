package main

import (
	"fmt"
	"pluginengine/utils"
)

// main takes:
// 1) IP of network device
// 2) function name (Discovery/Collect)
// 3) Metric Group Array(Only if 2nd parameter is collect
func main() {
	//var ip string = os.Args[1]
	fmt.Println(utils.Discovery("172.16.8.2"))

}
