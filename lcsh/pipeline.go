package lcsh

import (
	"io"

	"github.com/boutros/marc"
	"github.com/fubrenda/a/marctools"
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/recordstore"
)

func NewFileToRecordStorePipeline(rs *recordstore.RecordStore, data io.Reader, format marc.Format) *pipeline.Pipeline {
	marcReader := marctools.MustNewMarcReader(data, format)
	marcToReso := MustNewTransform(marcReader.Out)
	resoWriter := marctools.MustResoRecordWriter(rs, marcToReso.Out)

	return pipeline.MustNewPipeline("lcsh2db", marcReader, resoWriter, marcToReso)
}
