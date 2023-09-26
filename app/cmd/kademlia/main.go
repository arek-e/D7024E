package main

import (
	"fmt"
	"strconv"

	"github.com/arek-e/D7024E/app/cmd/cli"
	"github.com/arek-e/D7024E/app/internal"
	"github.com/arek-e/D7024E/app/utils"
)

var port = 1337

func main() {
	// Gets the docker containers IP
	localIP := utils.GetOutboundIP()
	fmt.Printf("LocalIP: %s\n", localIP.String())

	// Combines the ip with port 172.20.0.3 + ":" + port
	localAdress := fmt.Sprintf("%s:%d", localIP.String(), port)

	self := internal.NewKademliaNode(localAdress)

	network := &internal.Network{}
	network.Node = &self

	bootstrapNodeID := internal.NewRandomKademliaID()
	// Gets the boostrap ip address "172.20.0.2"
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

	cli := &cli.CLI{
		Node: &self,
		Net:  network,
	}
	// Start the CLI in a goroutine
	exitCh := make(chan struct{})
	go cli.StartCLI(exitCh)

	// Wait for the exit signal from the CLI
	<-exitCh
}
