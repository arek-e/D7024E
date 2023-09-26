package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type RPC struct {
	Type   string
	Sender Contact
	RpcID  *KademliaID
	Data   json.RawMessage
}

type PingRequest struct {
	PingID *KademliaID
}

type PingResponse struct {
	PongID *KademliaID
}

type FindContactRequest struct {
	Target *KademliaID
}

type FindContactResponse struct {
	Contacts []Contact
}

type StoreRequest struct {
	Key  string // Hashed key in the request
	Data string
}

type StoreResponse struct {
	KeyLocation string
}

type FindDataRequest struct {
	Hash string
}

type FindDataResponse struct {
	Data  []byte
	Nodes []Contact // Nodes that are close to the data
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
	case "PingRequest":
		var pingReq PingRequest
		if err := json.Unmarshal(request.Data, &pingReq); err != nil {
			log.Printf("Error unmarshaling PingRequest: %v", err)
			return RPC{}, err
		}
		MessageID := pingReq.PingID

		pingResponse := PingResponse{
			PongID: MessageID,
		}

		responseData, err := json.Marshal(pingResponse)
		if err != nil {
			log.Printf("Error marshaling PingResponse: %v", err)
			return RPC{}, err
		}

		response = RPC{
			Sender: network.Node.Self,
			Type:   "PingResponse",
			Data:   json.RawMessage(responseData),
			RpcID:  request.RpcID,
		}

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

	case "StoreRequest":
		var storeReq StoreRequest
		if err := json.Unmarshal(request.Data, &storeReq); err != nil {
			log.Printf("Error unmarshaling StoreRequest: %v", err)
			return RPC{}, err
		}

		network.Node.Datastore.putData(storeReq.Key, []byte(storeReq.Data))

		storeResponse := StoreResponse{
			KeyLocation: storeReq.Key,
		}

		responseData, err := json.Marshal(storeResponse)
		if err != nil {
			log.Printf("Error marshaling StoreResponse: %v", err)
			return RPC{}, err
		}

		response = RPC{
			Sender: network.Node.Self,
			Type:   "StoreResponse",
			Data:   json.RawMessage(responseData),
			RpcID:  request.RpcID,
		}

	case "FindDataRequest":
		var findDataReq FindDataRequest
		if err := json.Unmarshal(request.Data, &findDataReq); err != nil {
			log.Printf("Error unmarshaling FindDataRequest: %v", err)
			return RPC{}, err
		}
		var data []byte
		var foundHash bool
		data, foundHash = network.Node.getDataFromStore(findDataReq.Hash)

		if foundHash {

			findDataResponse := FindDataResponse{
				Data: data,
			}

			responseData, err := json.Marshal(findDataResponse)
			if err != nil {
				log.Printf("Error marshaling FindDataResponse: %v", err)
				return RPC{}, err
			}

			response = RPC{
				Sender: network.Node.Self,
				Type:   "FindDataResponse",
				Data:   json.RawMessage(responseData),
				RpcID:  request.RpcID,
			}

			return response, nil
		}

		// If the hash was not found then we get the contacts closer to the hash and return in order to update
		// shortlist
		contacts := network.Node.Routes.FindClosestContacts(NewKademliaID(findDataReq.Hash), 20)

		findDataResponse := FindDataResponse{
			Nodes: contacts,
		}

		responseData, err := json.Marshal(findDataResponse)
		if err != nil {
			log.Printf("Error marshaling FindDataResponse: %v", err)
			return RPC{}, err
		}

		response = RPC{
			Sender: network.Node.Self,
			Type:   "FindDataResponse",
			Data:   json.RawMessage(responseData),
			RpcID:  request.RpcID,
		}

	default:
		log.Printf("Unknown RPC request: %s", request.Type)
		return RPC{}, errors.New("Unknown RPC request type")
	}
	return response, nil
}

func (network *Network) ExtractResponseData(responseRPC RPC) (interface{}, error) {
	switch responseRPC.Type {
	case "PingResponse":
		var pingResponse PingResponse
		if err := json.Unmarshal(responseRPC.Data, &pingResponse); err != nil {
			return nil, err
		}
		return pingResponse, nil

	case "FindContactResponse":
		var findContactResponse FindContactResponse
		if err := json.Unmarshal(responseRPC.Data, &findContactResponse); err != nil {
			return nil, err
		}
		return findContactResponse, nil

	case "StoreResponse":
		var storeResponse StoreResponse
		if err := json.Unmarshal(responseRPC.Data, &storeResponse); err != nil {
			return nil, err
		}
		return storeResponse, nil

	case "FindDataResponse":
		var findDataResponse FindDataResponse
		if err := json.Unmarshal(responseRPC.Data, &findDataResponse); err != nil {
			return nil, err
		}
		return findDataResponse, nil

	default:
		return nil, fmt.Errorf("unknown Response Data type: %s", responseRPC.Type)
	}
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

	// Use a channel to signal when data is received or when the timeout occurs
	responseChan := make(chan RPC)
	errorChan := make(chan error)

	go func() {
		buf := make([]byte, 5000)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			errorChan <- fmt.Errorf("error reading data: %v", err)
			return
		}

		// Parse the received data to determine the response type
		parsedResponse, err := DeserializeRPC(buf[:n])
		if err != nil {
			errorChan <- fmt.Errorf("error parsing RPC response: %v", err)
			return
		}

		if Validate(request, parsedResponse) {
			network.Node.Routes.AddContact(parsedResponse.Sender)
		}

		responseChan <- parsedResponse
	}()

	// Use a select statement to wait for data or timeout
	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errorChan:
		return RPC{}, err
	case <-time.After(500 * time.Millisecond):
		network.Node.Routes.RemoveContact(*contact)
		return RPC{}, fmt.Errorf("timeout while waiting for UDP response")
	}
}
