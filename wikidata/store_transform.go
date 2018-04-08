package wikidata

import (
	"github.com/fubrenda/a/pipeline"
	"github.com/oliveagle/jsonpath"
	"github.com/rs/zerolog"
)

var ClaimsToDecode = []string{"P244"} // , "P214", "P4801", "P1014", "P486"

var PatternsToMatch = map[string]string{
	"Identifier":       "$.id",
	"Heading":          "$.labels[language,value]",
	"LCSHIdentifier":   "$.claims.P244[:].mainsnak.datavalue.value",
	"VIAFIdentifier":   "$.claims.P214[:].mainsnak.datavalue.value",
	"LCMARCIdentifier": "$.claims.P4801[:].mainsnak.datavalue.value",
	"AATIdentifier":    "$.claims.P1014[:].mainsnak.datavalue.value",
	"MESHIdentifier":   "$.claims.P486[:].mainsnak.datavalue.value",
}

type WikiRecord struct {
	Identifier       string            `json:"identifier,omitempty"`
	Heading          map[string]string `json:"heading,omitempty"`
	LCSHIdentifier   []string          `json:"lcsh_identifiers,omitempty"`
	VIAFIdentifier   []string          `json:"viaf_ddentifiers,omitempty"`
	LCMARCIdentifier []string          `json:"lcmarc_identifiers,omitempty"`
	AATIdentifier    []string          `json:"aat_identifiers,omitempty"`
	MESHIdentifier   []string          `json:"mesh_identifiers,omitempty"`
}

// WikiDataToWikirecordTransform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type WikiDataToWikirecordTransform struct {
	logger    zerolog.Logger
	In        chan []WikiDataEntity
	Out       chan []WikiRecord
	processed int64
	name      string
	patterns  map[string]*jsonpath.Compiled
}

// MustNewTransform creates a wiki data transform that filters
func MustNewWikiDataToWikirecordTransform(logger zerolog.Logger, in chan []WikiDataEntity) *WikiDataToWikirecordTransform {
	var patterns map[string]*jsonpath.Compiled
	for key, val := range PatternsToMatch {
		pat, err := jsonpath.Compile(val)
		if err != nil {
			panic(err)
		}
		patterns[key] = pat
	}
	return &WikiDataToWikirecordTransform{
		logger:    logger,
		In:        in,
		Out:       make(chan []WikiRecord),
		processed: 0,
		name:      "wikidata:wiki-to-internal",
		patterns:  patterns,
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

func Lookup(obj interface{}, pat *jsonpath.Compiled) interface{} {
	res, err := pat.Lookup(obj)
	if err != nil {
		return ""
	}

	return res
}

// Transform will filter records that don't have the keys we need
func (t *WikiDataToWikirecordTransform) Transform(messages []WikiDataEntity) []WikiRecord {
	msgs := make([]WikiRecord, 0)

	for _, msg := range messages {
		wikiRecord := WikiRecord{
			Identifier:       Lookup(msg, t.patterns["Identifier"]).(string),
			Heading:          Lookup(msg, t.patterns["Heading"]).(map[string]string),
			LCSHIdentifier:   Lookup(msg, t.patterns["LCSHIdentifier"]).([]string),
			VIAFIdentifier:   Lookup(msg, t.patterns["VIAFIdentifier"]).([]string),
			LCMARCIdentifier: Lookup(msg, t.patterns["LCMARCIdentifier"]).([]string),
			AATIdentifier:    Lookup(msg, t.patterns["AATIdentifier"]).([]string),
			MESHIdentifier:   Lookup(msg, t.patterns["MESHIdentifier"]).([]string),
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
