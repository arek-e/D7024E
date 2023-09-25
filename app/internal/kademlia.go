package internal

import (
	"sync"

	"github.com/arek-e/D7024E/app/utils"
)

type Kademlia struct {
	Self      Contact // NOTE: This might not be necessary since the routing table comes with "me"
	Routes    *RoutingTable
	Datastore *Datastore
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

func (kademlia *Kademlia) Lookup(targetOrHash interface{}) ([]Contact, []byte, Contact) {
	switch t := targetOrHash.(type) {
	case *KademliaID:
		// Handle contact lookup
		contacts := kademlia.LookupContact(t)
		return contacts, nil, Contact{}

	case string:
		// Handle data lookup
		data, con := kademlia.LookupData(t)
		return nil, data, con

	default:
		return nil, nil, Contact{}
	}
}

// LookupContact "...to locate the k closest nodes to some given node ID"
func (kademlia *Kademlia) LookupContact(target *KademliaID) (k_nodes []Contact) {
	network := &Network{}
	network.Node = kademlia
	ch := make(chan []Contact)  // Channel for response
	conCh := make(chan Contact) // Channel for response contact

	//"The first alpha contacts selected are used to create a shortlist for the search."
	shortlist := kademlia.NewLookupList(target)

	// If there are fewer than alpha contacts in that bucket, contacts are selected from other buckets.
	// The contact closest to the target key, closestNode, is noted.
	// min
	if shortlist.Len() < alpha {
		// If shortlist length is less than alpha, perform the lookup for the first node asynchronously.
		go PerformLookup(*target, shortlist.Nodes[0].Node, *network, ch, conCh)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < alpha; i++ {
			go PerformLookup(*target, shortlist.Nodes[i].Node, *network, ch, conCh)
		}
	}
	shortlist.updateLookupList(*target, ch, conCh, *network)

	// creating the result list
	for _, insItem := range shortlist.Nodes {
		k_nodes = append(k_nodes, insItem.Node)
	}

	return
}

func PerformLookup(targetID KademliaID, receiver Contact, net Network, ch chan []Contact, conCh chan Contact) {
	resultingNodes, _ := net.SendFindContactMessage(&receiver, &targetID)
	ch <- resultingNodes
	conCh <- receiver
}

// Given a hash from data, finds the closest node where the data is to be stored
func (kademlia *Kademlia) LookupData(hash string) ([]byte, Contact) {
	net := &Network{}
	net.Node = kademlia
	var waitgroup sync.WaitGroup

	hashID := NewKademliaID(hash)
	shortlist := kademlia.NewLookupList(hashID)

	ch := make(chan []Contact)          // channel -> returns contacts
	targetData := make(chan []byte)     // channel -> when the data is found it is communicated through this channel
	dataContactCh := make(chan Contact) // channel that only takes the contact that returned the data

	if shortlist.Len() < alpha {
		go PerformLookupData(hash, shortlist.Nodes[0].Node, *net, ch, targetData, dataContactCh)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < alpha; i++ {
			go PerformLookupData(hash, shortlist.Nodes[i].Node, *net, ch, targetData, dataContactCh)
		}
	}

	data, con := shortlist.updateLookupData(hash, ch, targetData, dataContactCh, *net, waitgroup)

	// creating the resultdata, con :=shortlist.updateLook list
	return data, con
}

func PerformLookupData(hash string, receiver Contact, net Network, ch chan []Contact, target chan []byte, dataContactCh chan Contact) {
	targetData, reslist, dataContact, _ := net.SendFindDataMessage(&receiver, hash)
	ch <- reslist
	target <- targetData
	dataContactCh <- dataContact
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
