package utils

//ips:
//1) 172.16.8.2
//2) 172.16.14.5
//3) 172.16.14.61

var ScalarOidToMetric = map[string]string{

	".1.3.6.1.2.1.1.5.0": "system.name",
	".1.3.6.1.2.1.1.1.0": "system.description",
	".1.3.6.1.2.1.1.6.0": "system.location",
	".1.3.6.1.2.1.1.2.0": "system.objectOID",
	".1.3.6.1.2.1.1.3.0": "system.uptime",
}
var MetricToScalarOid = map[string]string{

	"system.name":        ".1.3.6.1.2.1.1.5.0",
	"system.description": ".1.3.6.1.2.1.1.1.0",
	"system.location":    ".1.3.6.1.2.1.1.6.0",
	"system.objectOID":   ".1.3.6.1.2.1.1.2.0",
	"system.uptime":      ".1.3.6.1.2.1.1.3.0",
}

var InstanceOidToMetric = map[string]string{

	".1.3.6.1.2.1.2.2.1.1":     "interface.index",
	".1.3.6.1.2.1.31.1.1.1.1":  "interface.name",
	".1.3.6.1.2.1.2.2.1.8":     "interface.operational.status",
	".1.3.6.1.2.1.2.2.1.7":     "interface.admin.status",
	".1.3.6.1.2.1.31.1.1.1.18": "interface.alias",
	".1.3.6.1.2.1.2.2.1.2":     "interface.description",
	".1.3.6.1.2.1.2.2.1.20":    "interface.sent.error.packet",
	".1.3.6.1.2.1.2.2.1.14":    "interface.received.error.packet",
	".1.3.6.1.2.1.2.2.1.16":    "interface.sent.octets",
	".1.3.6.1.2.1.2.2.1.10":    "interface.received.octets",
	".1.3.6.1.2.1.2.2.1.5":     "interface.speed",
	".1.3.6.1.2.1.2.2.1.6":     "interface.physical.address",
}

var MetricToInstanceOid = map[string]string{

	"interface.index":                 ".1.3.6.1.2.1.2.2.1.1",
	"interface.name":                  ".1.3.6.1.2.1.31.1.1.1.1",
	"interface.operational.status":    ".1.3.6.1.2.1.2.2.1.8",
	"interface.admin.status":          ".1.3.6.1.2.1.2.2.1.7",
	"interface.alias":                 ".1.3.6.1.2.1.31.1.1.1.18",
	"interface.description":           ".1.3.6.1.2.1.2.2.1.2",
	"interface.sent.error.packet":     ".1.3.6.1.2.1.2.2.1.20",
	"interface.received.error.packet": ".1.3.6.1.2.1.2.2.1.14",
	"interface.sent.octets":           ".1.3.6.1.2.1.2.2.1.16",
	"interface.received.octets":       ".1.3.6.1.2.1.2.2.1.10",
	"interface.speed":                 ".1.3.6.1.2.1.2.2.1.5",
	"interface.physical.address":      ".1.3.6.1.2.1.2.2.1.6",
}
