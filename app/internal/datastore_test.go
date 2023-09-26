package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatastore(t *testing.T) {
	// Create a new Datastore
	datastore := NewDataStore()

	// Test adding data to the datastore
	key := "testKey"
	data := []byte("testData")
	datastore.putData(key, data)

	// Test getting data from the datastore
	retrievedData, exists := datastore.getData(key)
	assert.True(t, exists, "Expected data to exist in the datastore")
	assert.Equal(t, data, retrievedData, "Retrieved data should match the stored data")

	// Test getting non-existent data
	nonExistentKey := "nonExistentKey"
	_, exists = datastore.getData(nonExistentKey)
	assert.False(t, exists, "Expected non-existent key to not exist in the datastore")
}
