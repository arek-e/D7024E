package internal

import (
	"errors"
	"log"
	"sync"
	"time"
)

const TTL_AMOUNT = 10

type Datastore struct {
	Store map[string]*DataEntry
	TTL   time.Duration // U1.
}

type DataEntry struct {
	Data   []byte
	Time   time.Time  // U1.
	Forget bool       // U3.
	mu     sync.Mutex // Mutex for the specific entry
}

func NewDataStore() *Datastore {
	DS := &Datastore{}
	DS.Store = make(map[string]*DataEntry)
	DS.TTL = TTL_AMOUNT * time.Second

	return DS
}

func (DS *Datastore) putData(key string, data []byte) {
	entry := &DataEntry{
		Data:   data,
		Time:   DS.getExpirationTime(),
		Forget: false,
	}
	DS.Store[key] = entry
}

func (DS *Datastore) getData(key string) (val []byte, hasVal bool) {
	entry, found := DS.Store[key]
	if found {
		entry.mu.Lock()
		defer entry.mu.Unlock()

		if time.Now().After(entry.Time) {
			log.Printf("Data is expired: %v", key)
			delete(DS.Store, key)
			return nil, false
		}

		return entry.Data, true
	}
	return nil, false
}

func (DS *Datastore) getExpirationTime() (expirationTime time.Time) {
	return time.Now().Add(DS.TTL)
}

// U2.
func (DS *Datastore) refreshData(key string) error {
	entry, found := DS.Store[key]
	if !found {
		return errors.New("refreshData: key was not found")
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	entry.Time = DS.getExpirationTime()
	// No need to update DS.Store, as we're working with a reference
	return nil
}

// U3.
func (DS *Datastore) toggleForgetFlag(key string) error {
	log.Printf("Check hash %v", key)

	entry, found := DS.Store[key]
	if !found {
		return errors.New("toggleForgetFlag: key was not found")
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	entry.Forget = !entry.Forget
	// No need to update DS.Store, as we're working with a reference
	return nil
}

// U3.
func (DS *Datastore) checkForgetFlag(key string) bool {
	entry, found := DS.Store[key]
	if !found {
		return false
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	return entry.Forget
}
