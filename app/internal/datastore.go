package internal

// TODO: Add way to expire data
type Datastore struct {
	Store map[string][]byte
}

func NewDataStore() *Datastore {
	DS := &Datastore{}
	DS.Store = make(map[string][]byte)

	return DS
}

func (DS *Datastore) putData(key string, data []byte) {
	DS.Store[key] = data
}

func (DS *Datastore) getData(key string) (val []byte, hasVal bool) {
	val, hasVal = DS.Store[key]
	return
}
