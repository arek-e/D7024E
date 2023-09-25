package internal

type Kademlia struct {
	Self   Contact // NOTE: This might not be necessary since the routing table comes with "me"
	Routes RoutingTable
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
