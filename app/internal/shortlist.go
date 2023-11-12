package internal

import (
	"sort"
)

type ShortList struct {
	Nodes []ShortListItem
}

type ShortListItem struct {
	Node Contact
	Flag bool
}

// NewShortList returns a ShortList with k-closest nodes from the nodes routingtable.
func (kademlia *Kademlia) NewShortList(targetID *KademliaID) (shortlist *ShortList) {
	shortlist = &ShortList{}
	// "The first alpha (3) contacts selected are used to create a shortlist for the search. "
	closestK := kademlia.Routes.FindClosestContacts(targetID, alpha)

	for _, item := range closestK {
		lsItem := &ShortListItem{item, false}
		shortlist.Nodes = append(shortlist.Nodes, *lsItem)
	}
	return
}

func (shortlist *ShortList) refresh(contacts []Contact, notConsidered []ShortListItem) {
	candidateList := ShortList{}
	tempList := shortlist.Nodes

	// in the case that there were no contacts returned from the channels this would not run
	for _, contact := range contacts {
		listItem := ShortListItem{contact, false}
		candidateList.Nodes = append(candidateList.Nodes, listItem)
	}
	// Since the responder that sent 0 contacts has already been considered it is not a new candidate to consider
	sortingList := ShortList{}
	candidateList.Remove(notConsidered)
	// We add the nodes in the real shortlist to our temporary list to be sorted
	sortingList.Append(tempList)
	// Once again now if the already considered node was in the real shortlist it will be removed
	sortingList.Remove(notConsidered)
	// One final time we add all the new candidates that has not been considered (for loop)
	sortingList.Append(candidateList.Nodes)
	// We then sort based on the distance to the lookup node
	sortingList.Sort()

	// We overwrite the shortlist nodes with new candidates
	if len(sortingList.Nodes) < bucketSize {
		shortlist.Nodes = sortingList.GetContacts(len(sortingList.Nodes))
	} else {
		shortlist.Nodes = sortingList.GetContacts(bucketSize)
	}
}

func (shortlist *ShortList) updateShortList(targetID KademliaID, ch chan []Contact, conCh chan Contact, net Network) {
	consideredList := ShortList{}
	for {
		contacts := <-ch
		responder := <-conCh
		if len(contacts) > 0 {
			shortlist.refresh(contacts, consideredList.Nodes)
		} else {
			resItem := ShortListItem{responder, true}
			notConList := []ShortListItem{resItem}
			consideredList.Append(notConList)
			shortlist.refresh([]Contact{}, consideredList.Nodes)
		}
		// The refresh function has updated the shortlist with new canditates for lookup (or it is empty and we are done)
		nextContact, Done := shortlist.findNextLookup()
		if Done {
			return
		} else {
			go PerformLookup(targetID, nextContact, net, ch, conCh)
		}
	}
}

func (shortlist *ShortList) updateLookupData(hash string, ch chan []Contact, target chan []byte, dataContactCh chan Contact, net Network) ([]byte, Contact) {
	for {
		contacts := <-ch
		targetData := <-target
		dataContact := <-dataContactCh

		// data not nil = correct data is found
		if targetData != nil {
			return targetData, dataContact
		}

		shortlist.refresh(contacts, []ShortListItem{})
		nextContact, Done := shortlist.findNextLookup()
		if Done {
			return nil, Contact{}
		} else {
			go PerformLookupData(hash, nextContact, net, ch, target, dataContactCh)
		}
	}
}

func (shortlist *ShortList) findNextLookup() (Contact, bool) {
	var nextItem Contact
	done := true
	for i, item := range shortlist.Nodes {
		if !item.Flag {
			nextItem = item.Node
			shortlist.Nodes[i].Flag = true
			done = false
			break
		}
	}
	return nextItem, done
}

// Append an array of Contacts to the ContactCandidates if not duplicate
func (shortlist *ShortList) Append(Contacts []ShortListItem) {
	for _, newCandidate := range Contacts {
		add := true
		for _, candidate := range shortlist.Nodes {
			if candidate.Node.ID.Equals(newCandidate.Node.ID) {
				add = false
				break
			}
		}
		if add {
			shortlist.Nodes = append(shortlist.Nodes, newCandidate)
		}
	}
}

func (shortlist *ShortList) Remove(Contacts []ShortListItem) {
	for _, newCandidate := range Contacts {
		for i, candidate := range shortlist.Nodes {
			if candidate.Node.ID.Equals(newCandidate.Node.ID) {
				shortlist.remove(i)
				break
			}
		}
	}
}

// Len returns the lenght of the LookupList
func (shortlist *ShortList) Len() int {
	return len(shortlist.Nodes)
}

func (shortlist *ShortList) remove(n int) {
	shortlist.Nodes = append(shortlist.Nodes[:n], shortlist.Nodes[n+1:]...)
}

// GetContacts returns the first count number of Contacts
func (shortlist *ShortList) GetContacts(count int) []ShortListItem {
	return shortlist.Nodes[:count]
}

// Sort the Contacts in ContactCandidates
func (shortlist *ShortList) Sort() {
	sort.Sort(shortlist)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (shortlist *ShortList) Swap(i, j int) {
	shortlist.Nodes[i], shortlist.Nodes[j] = shortlist.Nodes[j], shortlist.Nodes[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (shortlist *ShortList) Less(i, j int) bool {
	return shortlist.Nodes[i].Node.Less(&shortlist.Nodes[j].Node)
}
