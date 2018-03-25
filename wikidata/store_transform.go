package wikidata

import (
	"encoding/json"
	"log"

	"github.com/fubrenda/a/pipeline"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

var ClaimsToDecode = []string{"P244"} // , "P214", "P4801", "P1014", "P486"

type WikiRecord struct {
	Identifier       string   `json:"identifier,omitempty"`
	Heading          []string `json:"heading,omitempty"`
	LCSHIdentifier   []string `json:"lcsh_identifiers,omitempty"`
	VIAFIdentifier   []string `json:"viaf_ddentifiers,omitempty"`
	LCMARCIdentifier []string `json:"lcmarc_identifiers,omitempty"`
	AATIdentifier    []string `json:"aat_identifiers,omitempty"`
	MESHIdentifier   []string `json:"mesh_identifiers,omitempty"`
}

func GetLocAuthorityID(claims *json.RawMessage) []string {
	value := make([]string, 0)

	rawData, err := claims.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	data := gjson.GetManyBytes(rawData, "P244.#.mainsnak.datavalue.value")
	for _, result := range data {
		value = append(value, result.String())
	}

	return value
}

func decodeClaims(message WikiDataMessage) WikiRecord {
	wikiRecord := WikiRecord{
		Identifier:     message.ID,
		LCSHIdentifier: GetLocAuthorityID(message.Claims),
	}

	return wikiRecord
}

// WikiDataToWikirecordTransform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type WikiDataToWikirecordTransform struct {
	logger    zerolog.Logger
	In        chan []WikiDataMessage
	Out       chan []WikiRecord
	processed int64
	name      string
}

// MustNewTransform creates a wiki data transform that filters
func MustNewWikiDataToWikirecordTransform(logger zerolog.Logger, in chan []WikiDataMessage) *WikiDataToWikirecordTransform {
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

// Transform will filter records that don't have the keys we need
func (t *WikiDataToWikirecordTransform) Transform(messages []WikiDataMessage) []WikiRecord {
	msgs := make([]WikiRecord, 0)

	for _, msg := range messages {
		msgs = append(msgs, decodeClaims(msg))
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
