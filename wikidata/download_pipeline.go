package wikidata

import (
	"io"

	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

func NewUrlToFilePipeline(logger zerolog.Logger, url string, fd io.Writer) *pipeline.Pipeline {
	messageReader := NewWikiDataReader(logger, url)
	messageFilter := MustNewTransform(logger, messageReader.Out)
	messageWriter := MustMessageWriter(logger, fd, messageFilter.Out)

	return pipeline.MustNewPipeline("wikidata2file", messageReader, messageWriter, messageFilter)
}
