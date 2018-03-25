package wikidata

import (
	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

type DBMessageWriter struct {
	logger  zerolog.Logger
	w       *WikiDataStore
	in      chan []WikiRecord
	name    string
	written int64
}

func MustDBMessageWriter(logger zerolog.Logger, w *WikiDataStore, in chan []WikiRecord) *DBMessageWriter {
	return &DBMessageWriter{
		logger:  logger,
		w:       w,
		in:      in,
		name:    "wikidata-db-writers",
		written: 0,
	}
}

func (w *DBMessageWriter) Write(killChan chan error) {
	for msgs := range w.in {
		err := w.w.SaveChunk(msgs)
		if err != nil {
			w.logger.Fatal().Err(err)
		}
		w.written = w.written + int64(len(msgs))
	}
}

func (w *DBMessageWriter) Name() string {
	return w.name
}

func (w *DBMessageWriter) Stats() pipeline.WriterStats {
	return pipeline.WriterStats{
		Written: w.written,
	}
}
