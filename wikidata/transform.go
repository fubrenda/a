package wikidata

import (
	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

// Transform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type Transform struct {
	logger       zerolog.Logger
	In           chan Message
	Out          chan Message
	processed    int64
	name         string
	expectedKeys []string
}

// MustNewTransform creates a wiki data transform that filters
func MustNewTransform(logger zerolog.Logger, in chan Message) *Transform {
	return &Transform{
		logger:       logger,
		In:           in,
		Out:          make(chan Message),
		processed:    0,
		name:         "wikidata:message-filter",
		expectedKeys: []string{"P244", "P214", "P4801", "P1014", "P486"},
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

// Transform will convert a chunk of marc.Records into a chunk of recordstore.ResoRecords
func (t *Transform) Transform(item Message) bool {
	claims := item["claims"].(map[string]interface{})
	for _, key := range t.expectedKeys {
		if _, ok := claims[key]; ok {
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
