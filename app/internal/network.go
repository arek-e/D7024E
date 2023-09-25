package internal

import (
	"log"
	"net"
	"strconv"

	"github.com/arek-e/D7024E/app/utils"
)

type Network struct {
	Node *Kademlia
}

func (network *Network) Listen(ip string, port int) {
	addr := utils.AddressToUDPAddr(ip + ":" + strconv.Itoa(port))

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Error listening on %s:%d: %v", addr.IP, addr.Port, err)
		return
	}
	defer conn.Close()

	log.Printf("Listening on: %s:%d", addr.IP, addr.Port)

	buffer := make([]byte, 1024)

	for {
		n, remoteaddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		receivedData := buffer[0:n]
		parsedRPCRequest, err := DeserializeRPC(receivedData)
		if err != nil {
			log.Printf("Error parsing RPC: %v", err)
			continue
		}

		network.Node.Routes.AddContact(parsedRPCRequest.Sender)
		responseRPC, err := network.CreateResponseRPC(parsedRPCRequest)
		if err != nil {
			log.Printf("Response error: %v", err)
			continue
		}

		serializedRPC, err := SerializeRPC(responseRPC)
		if err != nil {
			log.Printf("Response error: %v", err)
			continue
		}

		sendResponse(conn, remoteaddr, serializedRPC)
	}
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, responseMsg []byte) {
	_, err := conn.WriteToUDP([]byte(responseMsg), addr)
	if err != nil {
		log.Printf("Couldn't send response: %v", err)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
