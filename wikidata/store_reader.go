package wikidata

import (
	"encoding/json"
	"io"

	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

//WikiDataReader is a tool to read data from wikidata
type WikiDataFileReader struct {
	logger  zerolog.Logger
	chunker *MessageStream
	Out     chan []WikiDataEntity
	read    int64
	name    string
}

// NewWikiDataFileReader creates a new wiki data reader
func NewWikiDataFileReader(logger zerolog.Logger, reader io.Reader) (*WikiDataFileReader, error) {
	fileReader := &WikiDataFileReader{
		logger: logger,
		Out:    make(chan []WikiDataEntity, 10),
		chunker: &MessageStream{
			chunkSize: 1000,
			data:      json.NewDecoder(reader),
			done:      false,
		},
		name: "wikidata-file-reader",
	}

	return fileReader, nil
}

func (w *WikiDataFileReader) Read(killChan chan error) {
	w.logger.Info().Msg("Starting to read wikidata file")
	for w.chunker.Next() {
		w.Out <- w.chunker.Value()
		w.read += int64(len(w.chunker.Value()))
	}
}

// Finish close out the channel
func (w *WikiDataFileReader) Finish() {
	close(w.Out)
}

// Stats Report back stats on the reader
func (w *WikiDataFileReader) Stats() pipeline.ReaderStats {
	return pipeline.ReaderStats{
		Read: w.read,
	}
}

// Name returns the name of the reader
func (w *WikiDataFileReader) Name() string {
	return w.name
}
