package main

import (
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"os"
	"pluginengine/constants"
	"pluginengine/utils"
	"strconv"
	"strings"
)

// plugin
// Input json format
// {
// "id":string(discoveryId for discovery() and provisionId in case of polling()),
// "ip":string(device ip),
// "port":string,
// "community":string(default: public),
// "version":string(v1/v2c),
// "functionType":string(discovery/collect),
// "metricType":string(scalar/instance)(if functionType=collect)
// }

// Output json format
// {
// "id":string(reflect back what came in as id),
// "result": {
//			"status": "(success/failed)",
//			"message" : (plugin side message),
//			"system" : {}, (if scalar)
//			"interface":[{},{},{},...]  (if instance)
//			}
// }

func main() {

	var result = make(map[string]interface{})

	//at end, this function will print response object to stdout
	defer func() {

		if err := recover(); err != nil {

			//discovery profile id will be added on start itself
			//so no need to add here
			result[constants.RESULT] = nil

			result[constants.STATUS] = constants.FAILED

			result[constants.MESSAGE] = fmt.Sprintf("%v", err)

		}

		//TODO --> DONE takes 'id' which can be used for differing between multiple credential profiles
		//adding request context so that java(which has spawn this go exe)
		//can understand to which profile output belongs to in case of multiple IPs as command argument

		if strings.EqualFold(result[constants.STATUS].(string), constants.FAILED) {
			result[constants.RESULT] = nil
		}

		res, _ := json.Marshal(result)

		//giving plugin output here
		fmt.Println(string(res))

	}()

	var jsonArgs = make(map[string]interface{})

	var err = json.Unmarshal([]byte(os.Args[1]), &jsonArgs)

	if err != nil {

		panic(fmt.Sprintf("Error in Json formatting in Main(): %v", err))

	}

	//discovery id
	id, ok := jsonArgs["id"].(string)

	//adding so that it gets passed at end with result
	result["id"] = id

	if !ok {

		panic("cannot find Discovery ID")

	}

	ip, ok := jsonArgs["ip"].(string)

	if !ok {

		panic("cannot find IP")

	}

	//port will be given default 161
	port_, ok := jsonArgs["port"].(string)

	port, err := strconv.Atoi(port_)

	if err != nil || !ok {

		port = 161

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

	metricType, ok := jsonArgs["metricType"].(string)

	//only panics when collect functionType is given and no metric group is given
	if strings.EqualFold(functionType, "collect") && !ok {

		panic("cannot find metricType")

	}

	//setting configuration for snmp object
	gosnmp.Default.Target = ip

	gosnmp.Default.Port = uint16(port)

	gosnmp.Default.Community = community

	gosnmp.Default.Retries = 0

	switch {

	case strings.EqualFold(version, "v1"):

		gosnmp.Default.Version = gosnmp.Version1

	case strings.EqualFold(version, "v2c"):

		gosnmp.Default.Version = gosnmp.Version2c

	default:

		panic("Unsupported Version")

	}

	switch {

	//both collect and discovery returns full map with err(if any) or result

	//discovery will be same for v1 and v2c version as both will be using snmp get
	case strings.EqualFold(functionType, "discovery"):
		result = utils.Discovery(*gosnmp.Default)

	case strings.EqualFold(functionType, "collect"):
		result = utils.Collect(*gosnmp.Default, metricType)

	}

}
