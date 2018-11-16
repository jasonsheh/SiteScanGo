package info

import (
	"fmt"
	"net"
)

func GetPort(ipStr string) {
	ip := net.ParseIP(ipStr)
	portList := []int{80, 443, 3306, 8080}
	for _, port := range portList {
		tcpAddr := net.TCPAddr{
			IP:   ip,
			Port: port,
		}
		conn, err := net.DialTCP("tcp", nil, &tcpAddr)
		if err != nil {
			continue
		}
		conn.Close()
		fmt.Println(ip, port, "open")
	}
}
