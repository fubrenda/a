package lcsh

import (
	"io"

	"github.com/boutros/marc"
	"github.com/fubrenda/a/marctools"
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/recordstore"
	"github.com/fubrenda/a/wikidata"
)

func NewFileToRecordStorePipeline(rs *recordstore.RecordStore, wikidataStore *wikidata.WikiDataStore, data io.Reader, format marc.Format) *pipeline.Pipeline {
	marcReader := marctools.MustNewMarcReader(data, format)
	marcToReso := MustNewMarc2ResoTransform(marcReader.Out)
	enrichReso := MustNewEnrichResoTransform(marcToReso.Out, wikidataStore)
	resoWriter := marctools.MustResoRecordWriter(rs, enrichReso.Out)

	return pipeline.MustNewPipeline("lcsh2db", marcReader, resoWriter, marcToReso, enrichReso)
}
