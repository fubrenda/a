package wikidata

import (
	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
)

// WikiDataToWikirecordTransform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type WikiDataToWikirecordTransform struct {
	logger    zerolog.Logger
	In        chan []WikiDataEntity
	Out       chan []WikiRecord
	processed int64
	name      string
}

// MustNewTransform creates a wiki data transform that filters
func MustNewWikiDataToWikirecordTransform(logger zerolog.Logger, in chan []WikiDataEntity) *WikiDataToWikirecordTransform {

	return &WikiDataToWikirecordTransform{
		logger:    logger,
		In:        in,
		Out:       make(chan []WikiRecord),
		processed: 0,
		name:      "wikidata:wiki-to-internal",
	}
}

// Run will start the pipeline process
func (t *WikiDataToWikirecordTransform) Run(killChan chan error) {
	for item := range t.In {
		t.Out <- t.Transform(item)
		t.processed += int64(len(item))
	}
	close(t.Out)
}

func getExternalID(claimCode string, item WikiDataEntity) []string {
	ids := make([]string, 0)
	claims, ok := item.Claims[claimCode]
	if !ok {
		return ids
	}

	for _, claim := range claims {
		ids = append(ids, claim.Mainsnak.DataValue.Value.(string))
	}
	return ids
}

// Transform will filter records that don't have the keys we need
func (t *WikiDataToWikirecordTransform) Transform(messages []WikiDataEntity) []WikiRecord {
	msgs := make([]WikiRecord, 0)

	for _, msg := range messages {
		wikiRecord := WikiRecord{
			Identifier:       msg.ID,
			Heading:          msg.Labels,
			LCSHIdentifier:   getExternalID("P244", msg),
			VIAFIdentifier:   getExternalID("P214", msg),
			LCMARCIdentifier: getExternalID("P4801", msg),
			AATIdentifier:    getExternalID("P1014", msg),
			MESHIdentifier:   getExternalID("P486", msg),
		}

		msgs = append(msgs, wikiRecord)
	}

	return msgs
}

// Stats returns info about transform
func (t *WikiDataToWikirecordTransform) Stats() pipeline.TransformStats {
	return pipeline.TransformStats{
		Processed: t.processed,
	}
}

func (t *WikiDataToWikirecordTransform) Name() string {
	return t.name
}
