package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type RPC struct {
	Type   string
	Sender Contact
	RpcID  *KademliaID
	Data   json.RawMessage
}

type FindContactRequest struct {
	Target *KademliaID
}

type FindContactResponse struct {
	Contacts []Contact
}

func SerializeRPC(rpc RPC) ([]byte, error) {
	data, err := json.Marshal(rpc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DeserializeRPC(data []byte) (RPC, error) {
	var rpc RPC
	if err := json.Unmarshal(data, &rpc); err != nil {
		return RPC{}, err
	}

	return rpc, nil
}

func (network *Network) CreateResponseRPC(request RPC) (RPC, error) {
	var response RPC
	switch request.Type {

	case "FindContactRequest":
		var findContactReq FindContactRequest
		if err := json.Unmarshal(request.Data, &findContactReq); err != nil {
			log.Printf("Error unmarshaling FindContactRequest: %v", err)
			return RPC{}, err
		}
		target := findContactReq.Target
		contacts := network.Node.Routes.FindClosestContacts(target, bucketSize)

		findContactResponse := FindContactResponse{
			Contacts: contacts,
		}

		responseData, err := json.Marshal(findContactResponse)
		if err != nil {
			log.Printf("Error marshaling FindContactResponse: %v", err)
			return RPC{}, err
		}

		response = RPC{
			Sender: network.Node.Self,
			Type:   "FindContactResponse",
			Data:   json.RawMessage(responseData),
			RpcID:  request.RpcID,
		}

	default:
		log.Printf("Unknown RPC request: %s", request.Type)
		return RPC{}, errors.New("Unknown RPC request type")
	}
	return response, nil
}

func (network *Network) HandleResponseRPC(contact *Contact, request RPC) (RPC, error) {
	marshaledRPC, err := json.Marshal(request)
	if err != nil {
		return RPC{}, fmt.Errorf("error marshaling data: %v", err)
	}

	conn, err := network.sendRPC(contact, marshaledRPC)
	if err != nil {
		return RPC{}, fmt.Errorf("error sending UDP message: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 5000)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
	}

	// Parse the received data to determine the response type
	parsedResponse, err := DeserializeRPC(buf[:n])
	if err != nil {

	}

	return response, nil
}
