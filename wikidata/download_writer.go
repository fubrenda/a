package wikidata

import (
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

type MessageWriter struct {
	logger  zerolog.Logger
	writer  io.Writer
	in      chan Message
	name    string
	written int64
}

func MustMessageWriter(logger zerolog.Logger, w io.Writer, in chan Message) *MessageWriter {
	return &MessageWriter{
		logger:  logger,
		writer:  gzip.NewWriter(w),
		in:      in,
		name:    "wiki-writer",
		written: 0,
	}
}

func (w *MessageWriter) Write(killChan chan error) {
	for msg := range w.in {
		p, err := json.Marshal(msg)
		if err != nil {
			w.logger.Fatal().Err(err)
		}
		w.writer.Write(p)
		w.written++
	}
}

func (w *MessageWriter) Name() string {
	return w.name
}

func (w *MessageWriter) Stats() pipeline.WriterStats {
	return pipeline.WriterStats{
		Written: w.written,
	}
}
