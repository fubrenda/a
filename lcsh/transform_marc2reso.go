package lcsh

import (
	"github.com/boutros/marc"
	"github.com/fubrenda/a/marcdex"
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/recordstore"
)

// Marc2ResoTransform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type Marc2ResoTransform struct {
	In        chan []*marc.Record
	Out       chan []recordstore.ResoRecord
	processed int64
	name      string
}

// MustNewMarc2ResoTransform creates a lcsh Transform
func MustNewMarc2ResoTransform(in chan []*marc.Record) *Marc2ResoTransform {
	return &Marc2ResoTransform{
		In:        in,
		Out:       make(chan []recordstore.ResoRecord),
		processed: 0,
		name:      "lcsh:marc2reso",
	}
}

// Run will start the pipeline process
func (t *Marc2ResoTransform) Run(killChan chan error) {
	for item := range t.In {
		t.Out <- t.Transform(item)
	}
	close(t.Out)
}

// Transform will convert a chunk of marc.Records into a chunk of recordstore.ResoRecords
func (t *Marc2ResoTransform) Transform(chunk []*marc.Record) []recordstore.ResoRecord {
	output := make([]recordstore.ResoRecord, 0)
	for _, record := range chunk {
		output = append(output, marcdex.ConvertMarctoResoRecord(record))
		t.processed++
	}

	return output
}

// Stats returns info about transform
func (t *Marc2ResoTransform) Stats() pipeline.TransformStats {
	return pipeline.TransformStats{
		Processed: t.processed,
	}
}

func (t *Marc2ResoTransform) Name() string {
	return t.name
}
