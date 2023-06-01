package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"pluginengine/consts"
)

func Collect(ip string, metricType int) (result []byte, err any) {
	defer func() {
		if err = recover(); err != nil {
			log.Fatalf("Collect Function err: %v", err)
		}
	}()

	gosnmp.Default.Target = ip

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	err = gosnmp.Default.Connect()

	if err != nil {
		log.Fatalf("Collect Connect() err: %v", err)
		return
	}
	defer gosnmp.Default.Conn.Close()

	switch metricType {
	case 1:
		i := 0
		var oids = make([]string, len(consts.ScalarMetrics))
		for oid := range consts.ScalarMetrics {
			oids[i] = oid
			i++
		}
		res, err := gosnmp.Default.Get(oids)
		if err != nil {
			log.Fatalf("Collect Get() err: %v", err)
			return
		}

		result, err = json.Marshal(res)
		if err != nil {
			log.Fatalf("Collect Get() err: %v", err)
			return nil, err
		}

	case 2:
		for oid := range consts.InstanceMetrics {

			//callback while walk each value of instance metric
			var walkFunc = func(pdu gosnmp.SnmpPDU) error {
				fmt.Printf("%s = ", pdu.Name)
				switch pdu.Type {
				case gosnmp.OctetString:
					b := pdu.Value.([]byte)
					fmt.Printf("STRING: %s\n", string(b))
				default:
					fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
				}
				return nil
			}
			err = gosnmp.Default.BulkWalk(oid, walkFunc)

			if err != nil {
				log.Fatalf("Collect BulkWalk(%s) err: %v", oid, err)
			}
		}

	default:
		panic("Invalid Metric Group")

	}
	//TODO write function to transform result
	return (result), nil
}
