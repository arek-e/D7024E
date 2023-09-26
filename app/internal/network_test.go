package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	requestID := NewRandomKademliaID()
	requestRpc := RPC{RpcID: requestID}
	responseRpc := RPC{RpcID: requestID}

	result := Validate(requestRpc, responseRpc)

	assert.False(t, result, "Expected Validate to return false")

	// Test cases
	testCases := []struct {
		Request  RPC
		Response RPC
		Expected bool
	}{
		{
			Request:  RPC{Type: "PingRequest", RpcID: requestID},
			Response: RPC{Type: "PingResponse", RpcID: requestID},
			Expected: true,
		},
		{
			Request:  RPC{Type: "StoreRequest", RpcID: requestID},
			Response: RPC{Type: "StoreResponse", RpcID: requestID},
			Expected: true,
		},
		{
			Request:  RPC{Type: "FindContactRequest", RpcID: requestID},
			Response: RPC{Type: "FindContactResponse", RpcID: requestID},
			Expected: true,
		},
		{
			Request:  RPC{Type: "FindDataRequest", RpcID: requestID},
			Response: RPC{Type: "FindDataResponse", RpcID: requestID},
			Expected: true,
		},
		// Invalid cases
		{
			Request:  RPC{Type: "PingRequest", RpcID: requestID},
			Response: RPC{Type: "StoreResponse", RpcID: requestID},
			Expected: false,
		},
		{
			Request:  RPC{Type: "FindDataRequest", RpcID: requestID},
			Response: RPC{Type: "PingResponse", RpcID: requestID},
			Expected: false,
		},
		// Different requestID
		{
			Request:  RPC{Type: "PingRequest", RpcID: requestID},
			Response: RPC{Type: "PingResponse", RpcID: NewRandomKademliaID()},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run("Test Validate", func(t *testing.T) {
			result := Validate(testCase.Request, testCase.Response)
			assert.Equal(t, testCase.Expected, result)
		})
	}
}

func TestSendPingMessage(t *testing.T) {
	// Create a bootstrap node
	bootstrapAddress := "127.0.0.1:1351"
	bootstrapNode := NewKademliaNode(bootstrapAddress)

	// Create the second node
	secondNodeAddress := "127.0.0.1:1352"
	secondNode := NewKademliaNode(secondNodeAddress)

	// Create a simulated network for the bootstrap node
	bootstrapNetwork := &Network{}
	bootstrapNetwork.Node = &bootstrapNode

	// Start listening on the bootstrap node's address
	go bootstrapNetwork.Listen("127.0.0.1", 1351)

	// Create a simulated network for the second node
	secondNetwork := &Network{}
	secondNetwork.Node = &secondNode

	// Start listening on the second node's address
	go secondNetwork.Listen("127.0.0.1", 1352)

	// Perform the join operation for the second node
	_ = secondNode.JoinNetwork(&bootstrapNode.Self)

	// Create a contact for the second node
	contact := Contact{
		Address: secondNodeAddress,
	}

	// Call SendPingMessage on the network
	pingResponse, err := bootstrapNetwork.SendPingMessage(&contact)

	// Assert that there's no error
	assert.NoError(t, err, "Expected no error from SendPingMessage")

	// Assert that the Ping response is not nil
	assert.NotNil(t, pingResponse, "Expected non-nil PingResponse")
}
