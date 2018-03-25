package wikidata

import (
	"io"

	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

// TODO (ak): Rename store_pipeline stuff to index_* stuff

func NewFileToDBPipeline(logger zerolog.Logger, r io.Reader, store *WikiDataStore) *pipeline.Pipeline {
	messageReader, err := NewWikiDataFileReader(logger, r)
	if err != nil {
		logger.Error().Err(err)
	}
	messageTransform := MustNewWikiDataToWikirecordTransform(logger, messageReader.Out)
	messageWriter := MustDBMessageWriter(logger, store, messageTransform.Out)

	return pipeline.MustNewPipeline("wikidatafile2db", messageReader, messageWriter, messageTransform)
}
