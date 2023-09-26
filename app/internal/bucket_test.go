package internal

import (
	"github.com/arek-e/D7024E/app/utils"
	"testing"
)

func TestRemove(t *testing.T) {
	// create a new bucket
	bucket := newBucket()

	// Create a test contact and add to bucket
	target := NewContact(NewKademliaID(utils.Hash("localhost:1337")), "localhost:1337")
	bucket.AddContact(target)

	if bucket.Len() == 0 {
		t.Errorf("Fail: Contact was not added to the bucket.")
	}

	// Remove target from the bucket
	bucket.RemoveContact(target)

	if bucket.Len() > 0 {
		t.Errorf("Fail: Contact was not removed from the bucket.")
	}

}
