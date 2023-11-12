package internal

import (
	"log"
	"sync"
	"time"

	"github.com/arek-e/D7024E/app/utils"
)

type Kademlia struct {
	Self      Contact // NOTE: This might not be necessary since the routing table comes with "me"
	Routes    *RoutingTable
	Datastore *Datastore
	mu        sync.Mutex
}

// A system-wide concurrency parameter, such as 3.
const alpha int = 3

func NewKademliaNode(address string) (node Kademlia) {
	id := NewKademliaID(utils.Hash(address))
	node.Self = NewContact(id, address) // and store to contact object
	node.Routes = NewRoutingTable(node.Self)
	node.Datastore = NewDataStore()

	return
}

// JoinNetwork To join the network, a node u (self) must have a contact to an already participating node w (bootstrap). u inserts w into
// the appropriate k-bucket. u then performs a node lookup for its own node ID. Finally, u refreshes all k-
// buckets further away than its closest neighbor. During the refreshes, u both populates its own k-buckets
// and inserts itself into other nodes’ k-buckets as necessary
func (u *Kademlia) JoinNetwork(w *Contact) []Contact {
	// Add the bootstrap do the routing table
	u.Routes.AddContact(*w)
	// Perform a lookup on ourself
	u.mu.Lock()
	contacts, _, _ := u.Lookup(u.Self.ID)
	u.mu.Unlock()

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
	ch := make(chan []Contact)
	conCh := make(chan Contact)

	//"The first alpha contacts selected are used to create a shortlist for the search."
	shortlist := kademlia.NewShortList(target)

	// The contact closest to the target key, closestNode, is noted.
	if shortlist.Len() < alpha {
		// If shortlist length is less than alpha, perform the lookup for the first node.
		go PerformLookup(*target, shortlist.Nodes[0].Node, *network, ch, conCh)
	} else {
		//"The node then sends parallel, asynchronous FIND_* RPCs to the alpha contacts in the shortlist."
		for i := 0; i < alpha; i++ {
			go PerformLookup(*target, shortlist.Nodes[i].Node, *network, ch, conCh)
		}
	}
	shortlist.updateShortList(*target, ch, conCh, *network)

	// creating the result list
	for _, insItem := range shortlist.Nodes {
		k_nodes = append(k_nodes, insItem.Node)
	}

	return
}

// Given a hash from data, finds the closest node where the data is to be stored
func (kademlia *Kademlia) LookupData(hash string) ([]byte, Contact) {
	net := &Network{}
	net.Node = kademlia

	hashID := NewKademliaID(hash) // create kademlia ID from the hashed data
	shortlist := kademlia.NewShortList(hashID)

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

	data, con := shortlist.updateLookupData(hash, ch, targetData, dataContactCh, *net)

	// creating the resultdata, con :=shortlist.updateLook list
	return data, con
}

func PerformLookup(targetID KademliaID, receiver Contact, net Network, ch chan []Contact, conCh chan Contact) {
	resultingNodes, _ := net.SendFindContactMessage(&receiver, &targetID)
	ch <- resultingNodes
	conCh <- receiver
}

// runs SendFindDataMessage and loads response into two channels:
// ch -> contacts close to the data hash
// target -> the target data
func PerformLookupData(hash string, receiver Contact, net Network, ch chan []Contact, target chan []byte, dataContactCh chan Contact) {
	targetData, reslist, dataContact, _ := net.SendFindDataMessage(&receiver, hash)
	ch <- reslist
	target <- targetData
	dataContactCh <- dataContact
}

func (kademlia *Kademlia) Store(data []byte) (key string) {
	net := &Network{}
	net.Node = kademlia
	key = utils.Hash(string(data))

	kademlia.mu.Lock()
	kademlia.Datastore.putData(key, data)
	hashID := NewKademliaID(key)
	contactsToStore, _, _ := kademlia.Lookup(hashID)
	kademlia.mu.Unlock()

	for _, target := range contactsToStore {

		net.SendStoreMessage(data, &target)

		// U2.
		go func(contact Contact, key string) {
			refreshTicker := time.NewTicker(kademlia.Datastore.TTL / 2)
			defer refreshTicker.Stop()

			for {
				select {
				case <-refreshTicker.C:
					if kademlia.Datastore.checkForgetFlag(key) {
						return
					}

					refreshedContact, err := net.SendRefreshMessage(&contact, key)
					if err != nil {
						log.Printf("Error when refreshing: %v", err)
					}
					// TODO: Change so that a log file
					log.Printf("Refreshed data at: %v", refreshedContact.Address)

				case <-time.After(kademlia.Datastore.TTL):
					log.Println("TTL elapsed. Exiting goroutine.")
					return
				}
			}
		}(target, key)
	}

	return
}

func (kademlia *Kademlia) Refresh(hash string) (err error) {
	err = kademlia.Datastore.refreshData(hash)
	if err != nil {
		return err
	}
	return
}

// U3.
func (kademlia *Kademlia) Forget(hash string) (err error) {
	err = kademlia.Datastore.toggleForgetFlag(hash)
	if err != nil {
		return err
	}
	return
}

// GetDataFromStore(key) returns value and boolean
func (kademlia *Kademlia) getDataFromStore(hash string) (val []byte, hasVal bool) {
	val, hasVal = kademlia.Datastore.getData(hash)
	return
}
