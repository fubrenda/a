package wikidata

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

type messageStream struct {
	chunkSize      int
	scanner        *bufio.Scanner
	done           bool
	currentResults []Message
}

func (ms *messageStream) Next() bool {
	if ms.done {
		return false
	}
	ms.currentResults = make([]Message, 0)
	for ms.scanner.Scan() {
		var msg Message
		json.Unmarshal(ms.scanner.Bytes(), &msg)
		ms.currentResults = append(ms.currentResults, msg)
		if len(ms.currentResults) >= ms.chunkSize {
			return true
		}
	}

	ms.done = true

	return false
}

func (ms *messageStream) Value() []Message {
	return ms.currentResults
}

//WikiDataReader is a tool to read data from wikidata
type WikiDataFileReader struct {
	logger  zerolog.Logger
	chunker *messageStream
	Out     chan []Message
	read    int64
	name    string
}

// NewWikiDataFileReader creates a new wiki data reader
func NewWikiDataFileReader(logger zerolog.Logger, reader io.Reader) (*WikiDataFileReader, error) {
	zreader, err := gzip.NewReader(reader)
	if err != nil {
		logger.Error().Err(err)
		return nil, err
	}
	fileReader := &WikiDataFileReader{
		logger: logger,
		Out:    make(chan []Message, 10),
		chunker: &messageStream{
			chunkSize: 1000,
			scanner:   bufio.NewScanner(zreader),
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
		w.read = w.read + int64(len(w.chunker.Value()))
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
