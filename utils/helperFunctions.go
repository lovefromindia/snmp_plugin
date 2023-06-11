package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"pluginengine/constants"
	"strings"
)

// SnmpTypeConversion converts the snmp data types to go data types
func SnmpTypeConversion(pdu gosnmp.SnmpPDU) (result string) {

	switch pdu.Type {

	case gosnmp.OctetString:

		//checks for physical address oid prefix
		if strings.HasPrefix(pdu.Name, MetricToInstanceOid["interface.physical.address"]) {

			result = hex.EncodeToString(pdu.Value.([]byte))

		} else {

			result = fmt.Sprintf("%s", string(pdu.Value.([]byte)))

		}

	default:

		result = fmt.Sprintf("%v", pdu.Value)

	}

	return result

}

// GetDefaultResultMap returns default error map with status and err which is passed
func GetDefaultResultMap(status string, err error) map[string]interface{} {

	result := make(map[string]interface{})

	result[constants.STATUS] = status

	result[constants.MESSAGE] = fmt.Sprintf("%v", err)

	return result

}
