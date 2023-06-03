package utils

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"net"
	"pluginengine/consts"
)

// Discovery : this function will get scalar oid value to check if network device is responding
func Discovery(snmp gosnmp.GoSNMP) map[string]interface{} {
	result := make(map[string]interface{})

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	err := snmp.Connect()
	if err != nil {
		result["status"] = "failed"
		result["message"] = fmt.Errorf("connect() in Discovery function failed: %v", err)
		return result
	}

	defer func(Conn net.Conn) {
		tempErr := Conn.Close()
		if tempErr != nil {
			err = fmt.Errorf("close() in Discovery function failed: %v", tempErr)
		}
	}(snmp.Conn)

	//system.name oid to check if results are coming
	//to confirm that device is responding
	return getScalarOID(&snmp, []string{consts.MetricToScalarOid["system.name"]})
}
