package main

import (
	"encoding/hex"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"log"
)

func printValue(pdu g.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	fmt.Printf("Type: %T\n", pdu.Value)
	switch pdu.Type {
	case g.OctetString:
		b := pdu.Value.([]byte)
		fmt.Printf("STRING: %s\n", hex.EncodeToString(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, g.ToBigInt(pdu.Value))
	}
	return nil
}

func main() {
	g.Default.Target = "172.16.8.1"
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oid := "1.3.6.1.2.1.2.2.1.6.43"
	err2 := g.Default.BulkWalk(oid, printValue) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
}
