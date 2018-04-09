package main

import (
	"os"

	"github.com/fubrenda/a/logzer"
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/wikidata"
)

func main() {
	logger := logzer.MustNewLogzer("fetch_wikidata", false, os.Stdout)
	outFile, err := os.Create("data/out.json.gz")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create output file")
	}
	pl := wikidata.NewUrlToFilePipeline(
		logger,
		"https://dumps.wikimedia.org/wikidatawiki/entities/latest-all.json.bz2",
		outFile,
	)

	pipeline.RunReporter(pl)

}
