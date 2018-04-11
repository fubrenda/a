package wikidata_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/coreos/bbolt"
	"github.com/fubrenda/a/wikidata"
	"github.com/stretchr/testify/assert"
)

// tempfile returns a temporary file path.
func tempfile() string {
	f, _ := ioutil.TempFile("", "bolt-")
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func getDb() *bolt.DB {
	path := tempfile()
	defer os.RemoveAll(path)

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic("db open error: " + err.Error())
	}
	return db
}

func TestMustNewWikiDataStore(t *testing.T) {
	wikidatadb := getDb()

	wikidatastore := wikidata.MustNewWikiDataStore(wikidatadb)
	assert.IsType(t, &wikidata.WikiDataStore{}, wikidatastore)
}

func TestNewStorageOperation(t *testing.T) {
	storeOp := wikidata.NewStorageOperation("b", "c", "d", []byte("e"))
	assert.Equal(t, "b", storeOp.Bucket)
	assert.Equal(t, []byte("c:d"), storeOp.Key)
	assert.Equal(t, []byte("e"), storeOp.Value)
}

func TestConvertMessageToStorageOperations(t *testing.T) {
	wikidataRecord := wikidata.WikiRecord{
		Identifier: "myid",
		Heading: wikidata.LabelMap{
			"en": wikidata.LanguageValue{
				Language: "en",
				Value:    "bbb",
			},
		},
		LCSHIdentifier: []string{
			"n2345",
		},
	}

	storageOps, err := wikidata.ConvertMessageToStorageOperations(wikidataRecord)
	assert.Nil(t, err)
	assert.Len(t, storageOps, 2)
}

func TestHandleOperation(t *testing.T) {
	storageOp := wikidata.StorageOperation{
		Bucket: "testing",
		Key:    []byte("a:b"),
		Value:  []byte("bbb"),
	}

	wikidatadb := getDb()

	wikidatadb.Batch(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("testing"))
		wikidata.HandleOperation(tx, storageOp)
		return nil
	})

	wikidatadb.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("testing"))
		val := b.Get([]byte("a:b"))
		assert.Equal(t, []byte("bbb"), val)

		return nil
	})
}
