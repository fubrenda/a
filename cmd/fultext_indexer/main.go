package main

import (
	"time"

	"github.com/coreos/bbolt"
	"github.com/fubrenda/a/cli"
	"github.com/fubrenda/a/recordstore"
	"github.com/fubrenda/a/search"
)

func main() {
	args := cli.GetArgs()

	db, err := bolt.Open(args.Dbpath, 0666, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})
	defer db.Close()

	if err != nil {
		panic(err)
	}
	recordStore := recordstore.MustNewRecordStore(db)
	recordStream, err := recordstore.NewRecordStream(recordStore, 1000)
	if err != nil {
		panic(err)
	}
	searchIndex := search.MustNewIndex(args.IndexPath)

	for recordStream.Next() {
		println("yo")
		records := recordStream.Value()
		searchIndex.BatchIndex(records)
	}

}
