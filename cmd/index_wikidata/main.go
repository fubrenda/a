package main

import (
	"compress/gzip"
	"os"
	"time"

	"github.com/coreos/bbolt"
	"github.com/fubrenda/a/logzer"
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/wikidata"
)

func main() {
	logger := logzer.MustNewLogzer("index_wikidata", false, os.Stdout)
	inFile, err := os.Open("./data/out.json.gz")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open output file")
	}
	defer inFile.Close()

	decodedFile, err := gzip.NewReader(inFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create gzip reader")
	}

	db, err := bolt.Open("./data/wikidata.db", 0666, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open wikidata.db")
	}

	wikiStore := wikidata.MustNewWikiDataStore(db)

	pl := wikidata.NewFileToDBPipeline(
		logger,
		decodedFile,
		wikiStore,
	)

	pipeline.RunReporter(pl)

}
