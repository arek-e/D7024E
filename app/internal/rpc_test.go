package internal

import (
	"github.com/arek-e/D7024E/app/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateResponseRPCForUnknownRequest(t *testing.T) {
	// Create a Network instance for testing
	network := &Network{}

	// Create a request with an unknown type
	request := RPC{
		Type: "UnknownRequestType",
	}

	// Call the CreateResponseRPC function
	response, err := network.CreateResponseRPC(request)

	// Assert that an error is returned
	assert.Error(t, err, "Expected an error for unknown request type")

	// Assert that the response is empty
	assert.Equal(t, RPC{}, response, "Expected an empty response for unknown request type")
}

func TestRetrieveNonExistentData(t *testing.T) {
	// Start the bootstrap node (only listening, not joining)
	bootstrapAddress := "127.0.0.1:1310"
	bootstrapNode := NewKademliaNode(bootstrapAddress)

	// Start the second node and simulate it joining the network with the bootstrap node
	secondNodeAddress := "127.0.0.1:1311"
	secondNode := NewKademliaNode(secondNodeAddress)

	// Simulate the second node joining the bootNetwork with the bootstrap node
	bootNetwork := &Network{}
	bootNetwork.Node = &bootstrapNode

	go bootNetwork.Listen("127.0.0.1", 1310)

	joinNetwork := &Network{}
	joinNetwork.Node = &bootstrapNode

	// Perform the join operation and get the contacts
	_ = secondNode.JoinNetwork(&bootstrapNode.Self)
	go joinNetwork.Listen("127.0.0.1", 1311)

	dataToStore := "Lagrar saker för testning"
	hash := secondNode.Store([]byte(dataToStore))
	assert.NotEmpty(t, hash)

	lookupHash := utils.Hash("Hash som inte finns")

	// Simulate retrieving the stored data
	_, retrievedData, _ := secondNode.Lookup(lookupHash)

	// Assert that the retrieved data matches the stored data
	assert.Nil(t, retrievedData)
}

// TODO: Test timeout
func TestHandleResponseRPCWithTimeout(t *testing.T) {
	//// Create a bootstrap node
	//bootstrapAddress := "127.0.0.1:1300"
	//bootstrapNode := NewKademliaNode(bootstrapAddress)
	//
	//// Create the second node
	//secondNodeAddress := "127.0.0.1:1301"
	//secondNode := NewKademliaNode(secondNodeAddress)
	//
	//// Create a simulated network for the bootstrap node
	//bootstrapNetwork := &Network{}
	//bootstrapNetwork.Node = &bootstrapNode
	//
	//// Start listening on the bootstrap node's address
	//go bootstrapNetwork.Listen("127.0.0.1", 1300)
	//
	//// Create a simulated network for the second node
	//secondNetwork := &Network{}
	//secondNetwork.Node = &secondNode
	//
	//// Start listening on the second node's address
	//go secondNetwork.Listen("127.0.0.1", 1301)
	//
	//// Perform the join operation for the second node
	//_ = secondNode.JoinNetwork(&bootstrapNode.Self)
	//
	//// Create a contact for an address that is not in the network
	//nonExistentNodeAddress := "127.0.0.1:1201" // This address is not part of the network
	//
	//// Create a contact for the non-existent node
	//contact := Contact{
	//	Address: nonExistentNodeAddress,
	//}
	//
	//// Call SendPingMessage on the network
	//pingResponse, err := bootstrapNetwork.SendPingMessage(&contact)
	//
	//// Assert that there's a timeout error
	//assert.Error(t, err, "Expected a timeout error from SendPingMessage")
	//assert.Contains(t, err.Error(), "timeout while sending UDP message")
	//
	//// Assert that the Ping response is empty
	//assert.Equal(t, (*KademliaID)(nil), pingResponse, "Expected an empty PingResponse")
}
