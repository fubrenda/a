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

func TestSaveChunk(t *testing.T) {
	chunk := []wikidata.WikiRecord{
		wikidata.WikiRecord{
			Identifier: "myid1",
			Heading: wikidata.LabelMap{
				"en": wikidata.LanguageValue{
					Language: "en",
					Value:    "aaa",
				},
			},
			LCSHIdentifier: []string{
				"n2345",
			},
		},
		wikidata.WikiRecord{
			Identifier: "myid2",
			Heading: wikidata.LabelMap{
				"en": wikidata.LanguageValue{
					Language: "en",
					Value:    "bbb",
				},
			},
			LCSHIdentifier: []string{
				"n2346",
			},
		},
	}

	wikidatadb := getDb()
	wikidataStore := wikidata.MustNewWikiDataStore(wikidatadb)
	wikidataStore.SaveChunk(chunk)
	entry, err := wikidataStore.FindByIdentifier("myid2")
	assert.Nil(t, err)
	assert.NotNil(t, entry)
	results, err := wikidataStore.FindManyByPrefixIdentifier(
		wikidata.LCSHIdentifierPrefix,
		[]string{
			"n2345",
			"n2346",
		},
	)
	if err != nil {
		panic(err)
	}
	_, ok1 := results["n2345"]
	assert.True(t, ok1)
	_, ok2 := results["n2346"]
	assert.True(t, ok2)
}
