package internal

import "github.com/arek-e/D7024E/app/utils"

type Kademlia struct {
	Self   Contact // NOTE: This might not be necessary since the routing table comes with "me"
	Routes *RoutingTable
}

// A system-wide concurrency parameter, such as 3.
const alpha int = 3

func NewKademliaNode(address string) (node Kademlia) {
	id := NewKademliaID(utils.Hash(address))
	node.Self = NewContact(id, address) // and store to contact object
	node.Routes = NewRoutingTable(node.Self)

	return
}

// JoinNetwork To join the network, a node u (self) must have a contact to an already participating node w (bootstrap). u inserts w into
// the appropriate k-bucket. u then performs a node lookup for its own node ID. Finally, u refreshes all k-
// buckets further away than its closest neighbor. During the refreshes, u both populates its own k-buckets
// and inserts itself into other nodes’ k-buckets as necessary
func (u *Kademlia) JoinNetwork(w *Contact) []Contact {
	u.Routes.AddContact(*w)
	contacts, _, _ := u.Lookup(u.Self.ID)

	return contacts
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
