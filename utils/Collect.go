package utils

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"net"
	"pluginengine/consts"
	"strings"
	"sync"
)

func Collect(snmp *gosnmp.GoSNMP) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	err = snmp.Connect()

	if err != nil {
		return result, fmt.Errorf("connect() in Collect function failed: %v", err)
	}

	defer func(Conn net.Conn) {
		tempErr := Conn.Close()
		if tempErr != nil {
			err = fmt.Errorf("close() in Collect function failed: %v", tempErr)
		}
	}(snmp.Conn)

	//this channel will be used by both goroutines
	//to send data to this function to aggregate both
	//goroutine data
	var sharedResult = make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(sharedResult)
	}(&wg)

	//getting scalar oids
	go func(snmp *gosnmp.GoSNMP, wg *sync.WaitGroup) {
		defer wg.Done()
		scalarOIDS := make([]string, len(consts.ScalarMetrics))
		i := 0
		for key := range consts.ScalarMetrics {
			scalarOIDS[i] = key
			i++
		}
		res, err := snmp.Get(scalarOIDS)
		var scalarResult = make(map[string]interface{})
		if err != nil {
			log.Printf("get() in Collect function failed: %v", err)
			for key := range consts.ScalarMetrics {
				scalarResult[key] = nil
			}
		} else {
			mapOIDResult(res, scalarResult, consts.METRIC_TYPE_SCALAR)
		}

		//sending group type to easily identify at other end of channel
		sharedResult <- map[string]interface{}{"group": "scalar", "result": scalarResult}
	}(snmp, &wg)

	//getting instance oids
	go func(snmp *gosnmp.GoSNMP, wg *sync.WaitGroup) {
		defer wg.Done()

		var instanceResult = make([]map[string]string, 100)

		//for each of instance oid, we will call BulKWalk
		for rootOID := range consts.InstanceMetrics {
			err = snmp.BulkWalk(rootOID, func(dataUnit gosnmp.SnmpPDU) error {

			})
			if err != nil {
				log.Printf("bulkWalk() for OID: %s failed", rootOID)
			} else {

			}
		}

		//sending group type to easily identify at other end of channel
		sharedResult <- map[string]interface{}{"group": "instance", "result": instanceResult}
	}(snmp, &wg)

	//gathering results from channels polling different groups
	for res := range sharedResult {
		grpType := res.(map[string]interface{})["group"].(string)
		switch {
		case strings.EqualFold(grpType, "scalar"):
			for oid, val := range res.(map[string]interface{})["result"].(map[string]interface{}) {
				result[oid] = val
			}
		case strings.EqualFold(grpType, "instance"):

		default:
			log.Print("unknown group type received in Collect() function")
		}
	}

	return result, nil
}
