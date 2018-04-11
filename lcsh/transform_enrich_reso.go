package lcsh

import (
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/recordstore"
	"github.com/fubrenda/a/wikidata"
)

// EnrichResoTransform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type EnrichResoTransform struct {
	In         chan []recordstore.ResoRecord
	Out        chan []recordstore.ResoRecord
	processed  int64
	name       string
	wikidatadb *wikidata.WikiDataStore
}

// MustNewEnrichResoTransform creates a lcsh Transform
func MustNewEnrichResoTransform(in chan []recordstore.ResoRecord, wikidatadb *wikidata.WikiDataStore) *EnrichResoTransform {
	return &EnrichResoTransform{
		In:         in,
		Out:        make(chan []recordstore.ResoRecord),
		processed:  0,
		wikidatadb: wikidatadb,
		name:       "lcsh:enrich-reso",
	}
}

// Run will start the pipeline process
func (t *EnrichResoTransform) Run(killChan chan error) {
	for item := range t.In {
		t.Out <- t.Transform(item)
	}
	close(t.Out)
}

// Transform will convert a chunk of marc.Records into a chunk of recordstore.ResoRecords
func (t *EnrichResoTransform) Transform(chunk []recordstore.ResoRecord) []recordstore.ResoRecord {
	for i, record := range chunk {
		chunk[i] = record
		t.processed++
	}

	return chunk
}

// Stats returns info about transform
func (t *EnrichResoTransform) Stats() pipeline.TransformStats {
	return pipeline.TransformStats{
		Processed: t.processed,
	}
}

func (t *EnrichResoTransform) Name() string {
	return t.name
}
