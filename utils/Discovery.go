package utils

import (
	g "github.com/gosnmp/gosnmp"
	"log"
	"pluginengine/consts"
)

// Discovery this function will get scalar oid value to check if network device is responding
func Discovery(ip string) (status bool, err any) {
	status = false
	defer func() {
		if err = recover(); err != nil {
			log.Fatalf("Discovery Function err: %v", err)
		}
	}()

	g.Default.Target = ip

	//if ip address is reachable or not will not
	//be known until we start to send packets in UDP
	//so this line will be happily executed even if ip is not correct
	err = g.Default.Connect()
	if err != nil {
		log.Fatalf("Discovery Connect() err: %v", err)
		return false, err
	}
	defer g.Default.Conn.Close()

	_, err = g.Default.Get([]string{consts.ScalarMetrics["system.name"]})
	if err != nil {
		log.Fatalf("Discovery Get() err: %v", err)
		return false, err
	}
	return true, nil
}
