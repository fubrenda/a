package wikidata

import (
	"fmt"
	"log"

	"github.com/fubrenda/a/pipeline"
	"github.com/oliveagle/jsonpath"
	"github.com/rs/zerolog"
)

// Transform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type Transform struct {
	logger    zerolog.Logger
	In        chan WikiDataMessage
	Out       chan WikiDataMessage
	processed int64
	name      string
	filters   []*jsonpath.Compiled
}

// MustNewTransform creates a wiki data transform that filters
func MustNewTransform(logger zerolog.Logger, in chan WikiDataMessage) *Transform {
	patterns := make([]*jsonpath.Compiled, 0)
	for _, key := range []string{"P244", "P214", "P4801", "P1014", "P486"} {
		pat, err := jsonpath.Compile(fmt.Sprintf("$.claims.[?(@.%s)][0].id", key))
		if err != nil {
			log.Fatal(err)
		}
		patterns = append(patterns, pat)
	}

	return &Transform{
		logger:    logger,
		In:        in,
		Out:       make(chan WikiDataMessage),
		processed: 0,
		name:      "wikidata:message-filter",
		filters:   patterns,
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
func (t *Transform) Transform(item WikiDataMessage) bool {
	for _, pat := range t.filters {
		res, err := pat.Lookup(item)
		if err != nil {
			continue
		}
		val := res.(string)
		if val != "" {
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
