package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"os"
	"pluginengine/utils"
	"strconv"
	"strings"
)

func main() {

	var result = make(map[string]interface{})

	//at end, this function will print response object to stdout
	defer func() {
		if err := recover(); err != nil {
			result["result"] = nil
			result["error"] = fmt.Sprintf("%v", err)
		}

		//TODO:
		//adding request context so that java(which has spawn this go exe)
		//can understand to which profile output belongs to in case of multiple IPs

		fmt.Println(result)
		res, _ := json.Marshal(result)
		fmt.Println(string(res))

	}()

	var jsonArgs = make(map[string]interface{})
	var err = json.Unmarshal([]byte(os.Args[1]), &jsonArgs)
	if err != nil {
		panic(fmt.Sprintf("Error in Json formatting in Main(): %v", err))
	}

	ip, ok := jsonArgs["ip"].(string)
	if !ok {
		panic("cannot find IP")
	}

	community, ok := jsonArgs["community"].(string)
	if !ok {
		community = "public"
	}

	version, ok := jsonArgs["version"].(string)
	if !ok {
		version = "v2c"
	}

	functionType, ok := jsonArgs["functionType"].(string)
	if !ok {
		panic("cannot find functionType")
	}

	port_, ok := jsonArgs["port"].(string)
	port, err := strconv.Atoi(port_)
	if err != nil || !ok {
		port = 161
	}

	//setting configuration for snmp object
	gosnmp.Default.Target = ip
	gosnmp.Default.Community = community
	if strings.EqualFold(version, "v1") {
		gosnmp.Default.Version = gosnmp.Version1
	} else {
		gosnmp.Default.Version = gosnmp.Version2c
	}
	gosnmp.Default.Port = uint16(port)

	switch {
	case strings.EqualFold(functionType, "Discovery"):
		res, err := utils.Discovery(gosnmp.Default)
		if err != nil {
			panic(err)
		}
		result["result"] = res
	case strings.EqualFold(functionType, "Collect"):
		res, err := utils.Collect(gosnmp.Default)
		if err != nil {
			panic(err)
		}
		result["result"] = res
	default:
		panic("Unknown Function Type")
	}

}
