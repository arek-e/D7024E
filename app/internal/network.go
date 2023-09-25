package internal

import (
	"encoding/json"
	"fmt"
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

func (network *Network) sendRPC(contact *Contact, rpcData []byte) (*net.UDPConn, error) {
	host, port, err := net.SplitHostPort(contact.Address)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("Error parsing port: %v", err)
		return nil, err
	}

	nodeAddr := net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: parsedPort,
	}

	conn, err := net.DialUDP("udp", nil, &nodeAddr)
	if err != nil {
		errorMessage := fmt.Sprintf("Error creating UDP connection: %v", err)
		log.Print(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}

	_, err = conn.Write(rpcData)
	if err != nil {
		log.Printf("Error writing data: %v", err)
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func Validate(request RPC, response RPC) bool {
	if request.RpcID == nil || response.RpcID == nil {
		return false
	}

	if *request.RpcID != *response.RpcID {
		return false
	}
	switch request.Type {
	case "PingRequest":
		if response.Type == "PingResponse" {
			return true
		}
	case "FindContactRequest":
		if response.Type == "FindContactResponse" {
			return true
		}
	case "FindDataRequest":
		if response.Type == "FindDataResponse" {
			return true
		}
	}

	return false
}

func (network *Network) SendPingMessage(contact *Contact) (*KademliaID, error) {
	pingRequest := PingRequest{
		PingID: NewRandomKademliaID(),
	}

	log.Printf("PING: %v", pingRequest.PingID)

	marshalledData, err := json.Marshal(pingRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the data: %v", err)
	}

	requestRPC := RPC{
		Sender: network.Node.Self,
		Type:   "PingRequest",
		Data:   json.RawMessage(marshalledData),
		RpcID:  NewRandomKademliaID(),
	}

	response, err := network.HandleResponseRPC(contact, requestRPC)
	if err != nil {
		return nil, err
	}

	pingResponse, err := network.ExtractResponseData(response)
	if err != nil {
		return nil, err
	}

	pingResp, ok := pingResponse.(PingResponse)
	if !ok {
		return nil, fmt.Errorf("expected PingResponse, but got %T", pingResp)
	}

	log.Printf("PONG: %v", pingResp.PongID)

	return pingResp.PongID, nil
}

func (network *Network) SendFindContactMessage(contact *Contact, target *KademliaID) ([]Contact, error) {
	findContactReq := FindContactRequest{
		Target: target,
	}

	requestData, err := json.Marshal(findContactReq)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the data: %v", err)
	}

	requestRPC := RPC{
		Type:   "FindContactRequest",
		Sender: network.Node.Self,
		RpcID:  NewRandomKademliaID(),
		Data:   json.RawMessage(requestData),
	}

	response, err := network.HandleResponseRPC(contact, requestRPC)
	if err != nil {
		return nil, err
	}

	findContactResponse, err := network.ExtractResponseData(response)
	if err != nil {
		return nil, err
	}

	findContactResp, ok := findContactResponse.(FindContactResponse)
	if !ok {
		return nil, fmt.Errorf("expected FindContactResponse, but got %T", findContactResponse)
	}

	contacts := findContactResp.Contacts

	return contacts, nil
}

func (network *Network) SendFindDataMessage(contact *Contact, hash string) ([]byte, []Contact, Contact, error) {
	findDataReq := FindDataRequest{
		Hash: hash,
	}

	requestData, err := json.Marshal(findDataReq)
	if err != nil {
		return nil, nil, Contact{}, fmt.Errorf("unable to marshal the data: %v", err)
	}

	requestRPC := RPC{
		Type:   "FindDataRequest",
		Sender: network.Node.Self,
		RpcID:  NewRandomKademliaID(),
		Data:   json.RawMessage(requestData),
	}

	response, err := network.HandleResponseRPC(contact, requestRPC)
	if err != nil {
		log.Printf("Response is invalid: %+v", response)
		return nil, nil, Contact{}, err
	}

	findDataResponse, err := network.ExtractResponseData(response)
	if err != nil {
		return nil, nil, Contact{}, err
	}

	findDataResp, ok := findDataResponse.(FindDataResponse)
	if !ok {
		return nil, nil, Contact{}, fmt.Errorf("expected FindDataResponse, but got %T", findDataResponse)
	}

	retreivedData := findDataResp.Data

	return retreivedData, findDataResp.Nodes, response.Sender, nil
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
