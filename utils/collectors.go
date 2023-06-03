package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"pluginengine/consts"
	"strings"
	"sync"
)

func getScalarOID(snmp *gosnmp.GoSNMP, oids []string) map[string]interface{} {
	result := make(map[string]interface{})
	result["result"] = make(map[string]interface{})

	res, err := snmp.Get(oids)
	if err != nil {
		result["status"] = "failed"
		result["message"] = fmt.Sprintf("getScalarOID function failed: %v", err)
		result["result"] = nil
		return result
	}

	for _, val := range res.Variables {
		switch val.Type {
		case gosnmp.OctetString:
			result["result"].(map[string]interface{})[val.Name] = fmt.Sprintf("%s", string(val.Value.([]byte)))
		default:
			result["result"].(map[string]interface{})[val.Name] = fmt.Sprintf("%v", val.Value)
		}
	}
	result["status"] = "success"
	return result
}

// get info of single metric for all interfaces
func getInstanceOID(snmp *gosnmp.GoSNMP, rootOid string, result map[string]interface{}) {

	err := snmp.BulkWalk(rootOid, func(res gosnmp.SnmpPDU) error {
		tempArr := strings.Split(res.Name, ".")
		interfaceIndex := tempArr[len(tempArr)-1]
		_, ok := result["result"].(map[string]interface{})[interfaceIndex]
		if !ok {
			result["result"].(map[string]interface{})[interfaceIndex] = make(map[string]interface{})
		}
		switch res.Type {
		case gosnmp.OctetString:
			if strings.EqualFold(rootOid, consts.MetricToInstanceOid["interface.physical.address"]) {
				result["result"].(map[string]interface{})[interfaceIndex].(map[string]interface{})[consts.InstanceOidToMetric[rootOid]] = hex.EncodeToString(res.Value.([]byte))
			} else {
				result["result"].(map[string]interface{})[interfaceIndex].(map[string]interface{})[consts.InstanceOidToMetric[rootOid]] = string(res.Value.([]byte))
			}

		default:
			result["result"].(map[string]interface{})[interfaceIndex].(map[string]interface{})[consts.InstanceOidToMetric[rootOid]] = fmt.Sprintf("%v", res.Value)
		}
		return nil
	})
	if err != nil {
		result["status"] = "failed"
		result["message"] = fmt.Sprintf("getInstanceOID function failed: %v", err)
		result["result"] = nil
	}
	result["status"] = "success"
}

// get info of single metric for all interfaces
// but in this we get number of interfaces so that we don't
// make slices of size 0
func getNInstanceOID(snmp *gosnmp.GoSNMP, s *sync.WaitGroup, rootOid string, ifIndex int) map[string]interface{} {
	defer s.Done()
	result := make(map[string]interface{})
	result["result"] = make([]interface{}, ifIndex)

	countInterface := 0
	err := snmp.BulkWalk(rootOid, func(res gosnmp.SnmpPDU) error {
		switch res.Type {
		case gosnmp.OctetString:
			result["result"].([]interface{})[countInterface] = string(res.Value.([]byte))
		default:
			result["result"].([]interface{})[countInterface] = fmt.Sprintf("%d", gosnmp.ToBigInt(res.Value))
		}
		countInterface++
		return nil
	})
	if err != nil {
		result["status"] = "failed"
		result["message"] = fmt.Sprintf("getNInstanceOID function failed: %v", err)
		result["result"] = nil
		return result
	}
	result["status"] = "success"
	result["result"] = result["result"].([]interface{})[:countInterface]
	result["ifIndex"] = countInterface
	return result
}
