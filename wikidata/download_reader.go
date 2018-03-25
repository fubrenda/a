package wikidata

import (
	"compress/bzip2"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

//WikiDataReader is a tool to read data from wikidata
type WikiDataReader struct {
	logger zerolog.Logger
	URL    string
	Out    chan WikiDataMessage
	read   int64
	name   string
}

//NewHttpResponse starts and http request
func NewHttpResponse(logger zerolog.Logger, url string) io.Reader {
	resp, err := http.Get(url)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Trying to fetch: %s", url)
	}
	logger.Info().Str("url", url).Msg("Opening connection")
	return resp.Body
}

func NewStreamUncompress(r io.Reader) io.Reader {
	return bzip2.NewReader(r)
}

func (w *WikiDataReader) Decoder(data io.Reader) {
	w.logger.Info().Msg("Starting decode stream")
	dec := json.NewDecoder(data)

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		w.logger.Fatal().Err(err)
	}

	// while the array contains values
	for dec.More() {
		var m WikiDataMessage
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			w.logger.Fatal().Err(err)
		}
		w.Out <- m
		w.read++
	}

}

// NewWikiDataReader creates a new wiki data reader
func NewWikiDataReader(logger zerolog.Logger, url string) *WikiDataReader {
	return &WikiDataReader{
		logger: logger,
		Out:    make(chan WikiDataMessage, 10),
		URL:    url,
		name:   "wikidata-reader",
	}
}

func (w *WikiDataReader) Read(killChan chan error) {
	w.logger.Info().Msg("Starting to read wikidata")
	resp := NewHttpResponse(w.logger, w.URL)
	decomressed := NewStreamUncompress(resp)

	w.Decoder(decomressed)
}

// Finish close out the channel
func (w *WikiDataReader) Finish() {
	close(w.Out)
}

// Stats Report back stats on the reader
func (w *WikiDataReader) Stats() pipeline.ReaderStats {
	return pipeline.ReaderStats{
		Read: w.read,
	}
}

// Name returns the name of the reader
func (w *WikiDataReader) Name() string {
	return w.name
}
