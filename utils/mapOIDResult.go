package utils

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"pluginengine/consts"
)

// storeOIDResult : takes snmpPacket(result of snmp methods), result and metric type
// from snmpPacket it extract out the {oid:value} and transforms it into {oid_metric_name:value}
// and stores these pairs into result map passed by caller to be filled
func mapOIDResult(res *gosnmp.SnmpPacket, result map[string]interface{}, metricsType int) {
	for _, val := range res.Variables {
		switch val.Type {
		case gosnmp.OctetString:
			if metricsType == consts.METRIC_TYPE_SCALAR {
				result[consts.ScalarMetrics[val.Name]] = fmt.Sprintf("%s", string(val.Value.([]byte)))
			} else if metricsType == consts.METRIC_TYPE_INSTANCE {
				result[consts.InstanceMetrics[val.Name]] = fmt.Sprintf("%s", string(val.Value.([]byte)))
			}
		default:
			if metricsType == consts.METRIC_TYPE_SCALAR {
				result[consts.ScalarMetrics[val.Name]] = fmt.Sprintf("%d", gosnmp.ToBigInt(val.Value))
			} else if metricsType == consts.METRIC_TYPE_INSTANCE {
				result[consts.InstanceMetrics[val.Name]] = fmt.Sprintf("%d", gosnmp.ToBigInt(val.Value))
			}
		}
	}
}
