package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// From Brett @ NTC
// https://github.com/reiver/go-telnet/blob/master/conn.go // Line 35
func telnet() {

	conn, err := net.Dial("tcp", "172.20.10.75:2039")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(conn, "show version")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("%s", status)
}
