package discovery

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"net"
	"pluginengine/constants"
	"pluginengine/utils"
)

// Discovery : this function will get scalar oid value to check if network device is responding
// it will be same for both v1 and v2c
func Discovery(snmp gosnmp.GoSNMP) map[string]interface{} {

	result := make(map[string]interface{})

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	//error will be caught when we actually send snmp commands get,walk etc
	err := snmp.Connect()

	if err != nil {

		return utils.GetDefaultResultMap("failed", fmt.Errorf("error in Discovery(): %v", err))

	}

	defer func(Conn net.Conn) {

		tempErr := Conn.Close()

		if tempErr != nil {

			//TODO something about this close err
			//this is not big issue as at end of go plugin, every resource will be closed
			err = fmt.Errorf("close() in Discovery function failed: %v", tempErr)

		}

	}(snmp.Conn)

	result[constants.RESULT] = make(map[string]interface{})

	res, err := snmp.Get([]string{utils.MetricToScalarOid["system.name"]})

	if err != nil {

		return utils.GetDefaultResultMap(constants.FAILED, fmt.Errorf("getScalarOID function failed: %v", err))

	}

	for _, val := range res.Variables {

		result[constants.RESULT].(map[string]interface{})[utils.ScalarOidToMetric[val.Name]] = utils.SnmpTypeConversion(val)

	}

	result[constants.STATUS] = constants.SUCCESS

	return result

}
