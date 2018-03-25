package main

import (
	"compress/bzip2"
	"os"
	"time"

	"github.com/coreos/bbolt"
	"github.com/fubrenda/a/logzer"
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/wikidata"
)

func main() {
	logger := logzer.MustNewLogzer("index_wikidata", false, os.Stdout)
	inFile, err := os.Open("./data/out.json.bz2")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open output file")
	}
	defer inFile.Close()

	decodedFile := bzip2.NewReader(inFile)

	db, err := bolt.Open("./data/wikidata.db", 0666, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})

	if err != nil {
		logger.Fatal().Err(err)
	}

	wikiStore := wikidata.MustNewWikiDataStore(db)

	pl := wikidata.NewFileToDBPipeline(
		logger,
		decodedFile,
		wikiStore,
	)

	pipeline.RunReporter(pl)

}
