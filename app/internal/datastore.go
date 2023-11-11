package internal

import (
	"errors"
	"log"
	"sync"
	"time"
)

const TTL_AMOUNT = 10

type Datastore struct {
	Store map[string]DataEntry
	TTL   time.Duration // U1.
	mu    sync.Mutex    // U6.
}

// Refactor the entry since all three share the same hash string
type DataEntry struct {
	Data   []byte
	Time   time.Time // U1.
	Forget bool      // U3.
}

func NewDataStore() *Datastore {
	DS := &Datastore{}
	DS.Store = make(map[string]DataEntry)
	DS.TTL = TTL_AMOUNT * time.Second

	return DS
}

func (DS *Datastore) putData(key string, data []byte) {
	entry := DataEntry{
		Data:   data,
		Time:   DS.getExpirationTime(),
		Forget: false,
	}
	DS.Store[key] = entry
}

func (DS *Datastore) getData(key string) (val []byte, hasVal bool) {
	DS.mu.Lock()
	defer DS.mu.Unlock()

	entry, found := DS.Store[key]
	if found {
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
	expirationTime = time.Now().Add(DS.TTL)
	return
}

// U2.
func (DS *Datastore) refreshData(key string) error {
	DS.mu.Lock()
	defer DS.mu.Unlock()

	entry, found := DS.Store[key]
	if !found {
		return errors.New("refreshData: key was not found")
	}
	entry.Time = DS.getExpirationTime()
	DS.Store[key] = entry
	return nil
}

// U3.
func (DS *Datastore) toggleForgetFlag(key string) error {
	DS.mu.Lock()
	defer DS.mu.Unlock()

	log.Printf("Check hash %v", key)

	entry, found := DS.Store[key]
	if !found {
		return errors.New("toggleForgetFlag: key was not found")
	}
	entry.Forget = !entry.Forget
	DS.Store[key] = entry
	return nil
}

// U3.
func (DS *Datastore) checkForgetFlag(key string) bool {
	DS.mu.Lock()
	defer DS.mu.Unlock()

	entry := DS.Store[key]
	return entry.Forget
}
