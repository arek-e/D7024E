package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKademliaNode(t *testing.T) {
	address := "127.0.0.1:1337"
	node := NewKademliaNode(address)
	node2 := NewKademliaNode(address)

	// Check if both nodes are not nil
	assert.NotNil(t, node)
	assert.NotNil(t, node2)

	// Check if the Kademlia IDs are the same
	assert.Equal(t, node.Self.ID, node2.Self.ID)

	// Add assertions to validate the node's properties
	assert.NotNil(t, node)
	assert.NotNil(t, node.Self)
	assert.NotNil(t, node.Routes)
	assert.NotNil(t, node.Datastore)
}

func TestJoinNetworkAndLookup(t *testing.T) {
	// Start the bootstrap node (only listening, not joining)
	bootstrapAddress := "127.0.0.1:1337"
	bootstrapNode := NewKademliaNode(bootstrapAddress)

	// Start the second node and simulate it joining the network with the bootstrap node
	secondNodeAddress := "127.0.0.1:1338"
	secondNode := NewKademliaNode(secondNodeAddress)

	// Simulate the second node joining the bootNetwork with the bootstrap node
	bootNetwork := &Network{}
	bootNetwork.Node = &bootstrapNode

	go bootNetwork.Listen("127.0.0.1", 1337)

	joinNetwork := &Network{}
	joinNetwork.Node = &bootstrapNode

	// Perform the join operation and get the contacts
	contacts := secondNode.JoinNetwork(&bootstrapNode.Self)
	go joinNetwork.Listen("127.0.0.1", 1338)

	// Assert that the contacts slice has a length greater than 0
	assert.NotNil(t, contacts)
	assert.Len(t, contacts, 1) // Change the length to the expected value
}

func TestStoreData(t *testing.T) {
	// Start the bootstrap node (only listening, not joining)
	bootstrapAddress := "127.0.0.1:1120"
	bootstrapNode := NewKademliaNode(bootstrapAddress)

	// Start the second node and simulate it joining the network with the bootstrap node
	secondNodeAddress := "127.0.0.1:1121"
	secondNode := NewKademliaNode(secondNodeAddress)

	// Simulate the second node joining the bootNetwork with the bootstrap node
	bootNetwork := &Network{}
	bootNetwork.Node = &bootstrapNode

	go bootNetwork.Listen("127.0.0.1", 1120)

	joinNetwork := &Network{}
	joinNetwork.Node = &bootstrapNode

	// Perform the join operation and get the contacts
	_ = secondNode.JoinNetwork(&bootstrapNode.Self)
	go joinNetwork.Listen("127.0.0.1", 1121)

	dataToStore := "Lagrar saker för testning"
	hash := secondNode.Store([]byte(dataToStore))
	assert.NotEmpty(t, hash)

	// Simulate retrieving the stored data
	_, retrievedData, _ := secondNode.Lookup(hash)

	// Assert that the retrieved data matches the stored data
	assert.Equal(t, []byte(dataToStore), retrievedData)
}

func TestLookupWithInvalidInput(t *testing.T) {
	// Create a Kademlia node
	nodeAddress := "127.0.0.1:1340"
	kademliaNode := NewKademliaNode(nodeAddress)

	// Define an invalid input (neither *KademliaID nor string hash)
	invalidInput := 12345 // Replace with your desired invalid input

	// Call the Lookup method with the invalid input
	contacts, data, contact := kademliaNode.Lookup(invalidInput)

	// Assert that the returned values are as expected for an invalid input
	assert.Nil(t, contacts, "Contacts should be nil for invalid input")
	assert.Nil(t, data, "Data should be nil for invalid input")
	assert.Equal(t, Contact{}, contact, "Contact should be an empty Contact for invalid input")
}
