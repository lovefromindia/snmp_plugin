package utils

import (
	g "github.com/gosnmp/gosnmp"
	"log"
	"pluginengine/consts"
)

func Collect(ip string, metricType int) (result any, err any) {
	defer func() {
		if err = recover(); err != nil {
			log.Fatalf("Collect Function err: %v", err)
		}
	}()

	g.Default.Target = ip

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	err = g.Default.Connect()

	if err != nil {
		log.Fatalf("Collect Connect() err: %v", err)
		return nil, err
	}
	defer g.Default.Conn.Close()

	switch metricType {
	case 1:
		i := 0
		var oids = make([]string, len(consts.ScalarMetrics))
		for oid := range consts.ScalarMetrics {
			oids[i] = oid
			i++
		}
		result, err = g.Default.Get(oids)
		if err != nil {
			log.Fatalf("Collect Get() err: %v", err)
			return nil, err
		}
	case 2:
		for oid := range consts.InstanceMetrics {

			//TODO: make callback func for each
			var walkFunc = func() {}
			err = g.Default.BulkWalk(oid, walkFunc)

			if err != nil {
				log.Fatalf("Collect BulkWalk(%s) err: %v", oid, err)
			}
		}

	default:
		panic("Invalid Metric Group")

	}
	//TODO write function to transform result
	return transformResult(result), nil
}

func walkFunc()
