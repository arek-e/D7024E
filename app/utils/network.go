package utils

import (
	"log"
	"net"
	"strconv"
	"strings"
)

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetBootstrapAddress(ipString string, port string) string {
	stringList := strings.Split(ipString, ".")
	value := stringList[1]
	bootstrapNodeAddress := "172." + value + ".0.2:" + port
	return bootstrapNodeAddress
}

// Creates a UDPAddr from a contacts ip address.
func AddressToUDPAddr(address string) net.UDPAddr {
	addr, port, _ := net.SplitHostPort(address)
	netAddr := net.ParseIP(addr)
	intPort, _ := strconv.Atoi(port)
	receiver := net.UDPAddr{
		IP:   netAddr,
		Port: intPort,
	}
	return receiver
}
