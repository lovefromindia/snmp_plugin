package collect

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"net"
	"pluginengine/constants"
	"pluginengine/utils"
	"strings"
)

const (
	INTERFACE = "interface"
	SYSTEM    = "system"
	SCALAR    = "scalar"
	INSTANCE  = "instance"
)

// Collect : this will get all oid value
func Collect(snmp gosnmp.GoSNMP, metricType string) map[string]interface{} {

	result := make(map[string]interface{})

	//where oid results will be stored
	result[constants.RESULT] = make(map[string]interface{})

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	err := snmp.Connect()

	if err != nil {

		return utils.GetDefaultResultMap(constants.FAILED, fmt.Errorf("error in collect() method: %v", err))

	}

	defer func(Conn net.Conn) {

		tempErr := Conn.Close()

		if tempErr != nil {

			err = fmt.Errorf("close() in Collect function failed: %v", tempErr)

		}

	}(snmp.Conn)

	//collect data as per metric group
	switch {

	//both v1 and v2c can use snmp.get()
	case strings.EqualFold(metricType, SCALAR):

		//making map for sending system oid data
		result[constants.RESULT].(map[string]interface{})[SYSTEM] = make(map[string]interface{})

		scalarOIDS := make([]string, len(utils.ScalarOidToMetric))

		i := 0

		for oid := range utils.ScalarOidToMetric {

			scalarOIDS[i] = oid

			i++

		}

		data, err := snmp.Get(scalarOIDS)

		if err != nil {

			return utils.GetDefaultResultMap(constants.FAILED, fmt.Errorf("getScalarOID function failed: %v", err))

		}

		for _, val := range data.Variables {

			result[constants.RESULT].(map[string]interface{})[SYSTEM].(map[string]interface{})[utils.ScalarOidToMetric[val.Name]] = utils.SnmpTypeConversion(val)

		}

		result[constants.STATUS] = constants.SUCCESS

		return result

	case strings.EqualFold(metricType, INSTANCE):

		//making map for storing interface oids data
		//with interface index as keys
		tempMap := make(map[string]interface{})

		errors := make([]string, 0)

		//walkOrBulkWalk: this variable is function pointer
		//to either BulkWalk or Walk based on snmp version
		var walkOrBulkWalk = snmp.BulkWalk

		if snmp.Version == gosnmp.Version1 {

			walkOrBulkWalk = snmp.Walk

		} else if snmp.Version == gosnmp.Version2c {

			walkOrBulkWalk = snmp.BulkWalk

		} else {
			return utils.GetDefaultResultMap("failed", fmt.Errorf("unsupported Snmp Version"))
		}

		for rootOid := range utils.InstanceOidToMetric {

			err = walkOrBulkWalk(rootOid, func(pdu gosnmp.SnmpPDU) error {

				tempArr := strings.Split(pdu.Name, ".")

				interfaceIndex := tempArr[len(tempArr)-1]

				//if not any data is inserted for that interface
				//then make map to store data. If we don't do this
				//and directly try to access by index we get nil map
				_, ok := tempMap[interfaceIndex]

				if !ok {

					tempMap[interfaceIndex] = make(map[string]interface{})

				}

				tempMap[interfaceIndex].(map[string]interface{})[utils.InstanceOidToMetric[rootOid]] = utils.SnmpTypeConversion(pdu)

				return nil

			})

			//store err for single rootOid(since we are in loop)
			if err != nil {

				errors = append(errors, fmt.Sprintf("%v", err))

			}

		}

		//if all rootOid fetch have some errors then its failed
		if len(errors) >= len(utils.InstanceOidToMetric) {

			result[constants.STATUS] = constants.FAILED

		} else {

			result[constants.STATUS] = constants.SUCCESS

		}

		//store errors (if any)
		result[constants.MESSAGE] = strings.Join(errors, "\n") //join() converts array of errors into string
		// as java side "message" is string and not json array

		result[constants.RESULT].(map[string]interface{})[INTERFACE] = make([]interface{}, len(tempMap))

		i := 0

		for _, data := range tempMap {

			result[constants.RESULT].(map[string]interface{})[INTERFACE].([]interface{})[i] = data

			i++

		}

	default:
		return utils.GetDefaultResultMap(constants.FAILED, fmt.Errorf("unknown metricType"))

	}

	return result

}
