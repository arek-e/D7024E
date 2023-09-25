package main

import (
	"fmt"
	"strconv"

	"github.com/arek-e/D7024E/app/internal"
	"github.com/arek-e/D7024E/app/utils"
)

var port = 1337

func main() {
	localIP := utils.GetOutboundIP()
	fmt.Printf("LocalIP: %s\n", localIP.String())

	localAdress := fmt.Sprintf("%s:%d", localIP.String(), port)

	self := internal.NewKademliaNode(localAdress)

	network := &internal.Network{}
	network.Node = &self

	bootstrapNodeID := internal.NewRandomKademliaID()
	bootstrapNodeAddress := utils.GetBootstrapAddress(localIP.String(), strconv.Itoa(port))
	bootstrapNodeContact := internal.NewContact(bootstrapNodeID, bootstrapNodeAddress)

	// checkar ifall noden finns eller inte nätverket. Om den inte gör så den med
	// checkar även ifall det är självaste bootstrap noden
	if localAdress != bootstrapNodeAddress {
		self.JoinNetwork(&bootstrapNodeContact)
	} else {
		fmt.Printf("Bootstrap node started listening\n")
	}

	go network.Listen(localIP.String(), port)
}
