package wikidata

import (
	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

// Transform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type Transform struct {
	logger             zerolog.Logger
	In                 chan WikiDataEntity
	Out                chan WikiDataEntity
	processed          int64
	name               string
	requiredClaimCodes []string
}

// MustNewTransform creates a wiki data transform that filters
func MustNewTransform(logger zerolog.Logger, in chan WikiDataEntity) *Transform {

	return &Transform{
		logger:             logger,
		In:                 in,
		Out:                make(chan WikiDataEntity),
		processed:          0,
		name:               "wikidata:message-filter",
		requiredClaimCodes: []string{"P244", "P214", "P4801", "P1014", "P486"},
	}
}

// Run will start the pipeline process
func (t *Transform) Run(killChan chan error) {
	for item := range t.In {
		if t.Transform(item) {
			t.Out <- item
		}
		t.processed++
	}
	close(t.Out)
}

// Transform will filter records that don't have the keys we need
func (t *Transform) Transform(item WikiDataEntity) bool {

	for _, code := range t.requiredClaimCodes {
		if _, ok := item.Claims[code]; ok {
			return true
		}
	}

	return false
}

// Stats returns info about transform
func (t *Transform) Stats() pipeline.TransformStats {
	return pipeline.TransformStats{
		Processed: t.processed,
	}
}

func (t *Transform) Name() string {
	return t.name
}
