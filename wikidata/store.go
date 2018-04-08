package wikidata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/coreos/bbolt"
	"github.com/vmihailenco/msgpack"
)

const (
	// WikiDataEntityBucketName is the name of DB Bucket
	WikiDataEntityBucketName = "WikiDataEntity"
	// IdentifierKeyPrefix prefixes keys in boltdb for identifier field
	IdentifierKeyPrefix = "identifier"
)

// RecordStore will store a record into an index
type WikiDataStore struct {
	db *bolt.DB
}

// StorageOperation is an operation that will get stored in a key value database
type StorageOperation struct {
	Key    []byte
	Value  []byte
	Bucket string
}

// NewStorageOperation crates a new KeyValue
func NewStorageOperation(bucket string, keyPrefix string, key string, value []byte) StorageOperation {
	return StorageOperation{
		Key:    []byte(fmt.Sprintf("%s:%s", keyPrefix, key)),
		Value:  value,
		Bucket: bucket,
	}
}

// MustNewRecordStore will create a new RecordStore
func MustNewWikiDataStore(db *bolt.DB) *WikiDataStore {
	err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(WikiDataEntityBucketName))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return &WikiDataStore{
		db: db,
	}
}

// ConvertMessageToStorageOperations returns a list of operations to be stored
func ConvertMessageToStorageOperations(message WikiRecord) ([]StorageOperation, error) {
	keyValues := make([]StorageOperation, 0)
	mainValue, err := msgpack.Marshal(message)

	if err != nil {
		return keyValues, nil
	}

	keyValues = append(keyValues, NewStorageOperation(WikiDataEntityBucketName, IdentifierKeyPrefix, message.Identifier, mainValue))

	return keyValues, nil
}

//HandleOperation will persist an operation into the database
func HandleOperation(tx *bolt.Tx, operation StorageOperation) error {
	bucket := tx.Bucket([]byte(operation.Bucket))
	val := bucket.Get(operation.Key)
	if val != nil {
		log.Printf("While fetching key %s found exisisting %s", string(operation.Key), string(operation.Value))
	}
	err := bucket.Put(operation.Key, operation.Value)

	if err != nil {
		return err
	}

	return nil
}

// SaveChunk persists a chunk of ResoRecords to database
func (r *WikiDataStore) SaveChunk(messages []WikiRecord) error {
	return r.db.Update(func(tx *bolt.Tx) error {

		for _, message := range messages {
			operations, err := ConvertMessageToStorageOperations(message)
			if err != nil {
				return fmt.Errorf("SaveChunk::ConvertRecordToKeyValues: %s", err)
			}
			for _, operation := range operations {
				err := HandleOperation(tx, operation)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// FindByIdentifier will lookup a ResoRecord by its main identifier
func (r *WikiDataStore) FindByIdentifier(id string) (*WikiDataEntity, error) {
	var message WikiDataEntity
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(WikiDataEntityBucketName))
		value := bucket.Get([]byte(fmt.Sprintf("%s:%s", IdentifierKeyPrefix, id)))

		msgpack.Unmarshal(value, &message)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &message, nil
}

type RecordPage struct {
	Prefix  []byte
	Records []WikiDataEntity
	LastKey []byte
	More    bool
}

func (r *WikiDataStore) Scan(prefix []byte, currentKey []byte, numResults int) RecordPage {
	records := make([]WikiDataEntity, 0)
	startingPrefix := prefix
	if bytes.Compare(currentKey, startingPrefix) == 1 {
		startingPrefix = currentKey
	}

	var recordPage RecordPage
	r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(WikiDataEntityBucketName))
		c := bucket.Cursor()
		results := 0
		for k, v := c.Seek(startingPrefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			results++
			if results > numResults {
				recordPage = RecordPage{
					Prefix:  prefix,
					Records: records,
					LastKey: []byte(k),
					More:    true,
				}
			}
			var message WikiDataEntity
			err := msgpack.Unmarshal(v, &message)
			if err != nil {
				return err
			}
			records = append(records, message)
		}

		recordPage = RecordPage{
			Prefix:  prefix,
			Records: records,
			LastKey: []byte(""),
			More:    false,
		}

		return nil
	})

	return recordPage
}

func (r *WikiDataStore) Stats() (string, error) {
	var statStr string
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(WikiDataEntityBucketName))
		stats := bucket.Stats()
		statBytes, err := json.Marshal(stats)
		if err != nil {
			return err
		}
		statStr = string(statBytes)
		return nil
	})

	if err != nil {
		return "", err
	}

	return statStr, nil
}
