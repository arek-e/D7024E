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

// NewLookupList retuns a LookupList with k-closest nodes from the nodes routingtable.
func (kademlia *Kademlia) NewLookupList(targetID *KademliaID) (shortlist *ShortList) {
	shortlist = &ShortList{}
	closestK := kademlia.Routes.FindClosestContacts(targetID, bucketSize)

	for _, item := range closestK {
		lsItem := &ShortListItem{item, false}
		shortlist.Nodes = append(shortlist.Nodes, *lsItem)
	}
	return
}

func (shortlist *ShortList) refresh(contacts []Contact, notConsidered []ShortListItem) {
	candidateList := ShortList{} // holds the response []Contact
	tempList := shortlist.Nodes  // Copy of lookuplist
	for _, contact := range contacts {
		listItem := ShortListItem{contact, false}
		candidateList.Nodes = append(candidateList.Nodes, listItem)
	}
	sortingList := ShortList{}

	candidateList.Remove(notConsidered)

	sortingList.Append(tempList)

	sortingList.Remove(notConsidered)

	sortingList.Append(candidateList.Nodes)

	sortingList.Sort()

	if len(sortingList.Nodes) < bucketSize {
		shortlist.Nodes = sortingList.GetContacts(len(sortingList.Nodes))
	} else {
		shortlist.Nodes = sortingList.GetContacts(bucketSize)
	}
}

func (lookuplist *ShortList) updateLookupList(targetID KademliaID, ch chan []Contact, conCh chan Contact, net Network) {
	notConsidered := ShortList{}
	for {
		contacts := <-ch
		responder := <-conCh
		if len(contacts) > 0 {
			lookuplist.refresh(contacts, notConsidered.Nodes)
		} else {
			resItem := ShortListItem{responder, true}
			notConList := []ShortListItem{resItem}
			notConsidered.Append(notConList)
			lookuplist.refresh([]Contact{}, notConsidered.Nodes)
		}
		nextContact, Done := lookuplist.findNextLookup()
		if Done {
			return
		} else {
			go PerformLookup(targetID, nextContact, net, ch, conCh)
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
