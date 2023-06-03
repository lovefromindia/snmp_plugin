package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"net"
	"pluginengine/consts"
	"strings"
	"sync"
)

// Collect : this will get all oids value
func Collect(snmp gosnmp.GoSNMP) map[string]interface{} {
	result := make(map[string]interface{})

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	err := snmp.Connect()

	if err != nil {
		result["status"] = "failed"
		result["message"] = fmt.Errorf("connect() in Collect function failed: %v", err)
		return result
	}

	defer func(Conn net.Conn) {
		tempErr := Conn.Close()
		if tempErr != nil {
			err = fmt.Errorf("close() in Collect function failed: %v", tempErr)
		}
	}(snmp.Conn)

	var wg sync.WaitGroup
	wg.Add(2)

	//getting scalar oids
	var scalarResult = make(map[string]interface{})
	go func(snmp *gosnmp.GoSNMP, wg *sync.WaitGroup) {
		defer wg.Done()
		scalarOIDS := make([]string, len(consts.ScalarOidToMetric))
		i := 0
		for oid := range consts.ScalarOidToMetric {
			scalarOIDS[i] = oid
			i++
		}
		tempResult := getScalarOID(snmp, scalarOIDS)

		//storing results in map for scalar metrics
		for oid, val := range tempResult["result"].(map[string]interface{}) {
			scalarResult[consts.ScalarOidToMetric[oid]] = val
		}
		scalarResult["status"] = "success"

	}(&snmp, &wg)

	//getting instance oids
	var instanceResult = make(map[string]interface{})
	instanceResult["result"] = make(map[string]interface{})
	go func(snmp *gosnmp.GoSNMP, wg *sync.WaitGroup) {
		defer wg.Done()
		for rootOid := range consts.InstanceOidToMetric {
			getInstanceOID(snmp, rootOid, instanceResult)
		}
	}(&snmp, &wg)

	wg.Wait()
	//filling values of both scalarResult and instanceResult into main result map
	result["status"] = "success"
	result["result"] = make(map[string]interface{})
	if strings.EqualFold(instanceResult["status"].(string), "success") {
		result["result"].(map[string]interface{})["interfaces"] = make([]interface{}, 0)
		for _, val := range instanceResult["result"].(map[string]interface{}) {
			result["result"].(map[string]interface{})["interfaces"] = append(result["result"].(map[string]interface{})["interfaces"].([]interface{}), val)
		}
	}

	if strings.EqualFold(scalarResult["status"].(string), "success") {
		for oid, val := range scalarResult {
			result["result"].(map[string]interface{})[oid] = val
		}
	}

	val, err := json.Marshal(result)
	fmt.Println(string(val))
	return result
}
