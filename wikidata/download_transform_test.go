package wikidata

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func buildWikiDataEntity(claimCode string) WikiDataEntity {
	return WikiDataEntity{
		ID:   "Q26",
		Type: "item",
		Labels: map[string]LanguageValue{
			"en-gb": LanguageValue{
				Language: "en-gb",
				Value:    "Northern Ireland",
			},
		},
		Aliases: map[string][]LanguageValue{
			"ca": []LanguageValue{CaUlster},
		},
		SiteLinks: map[string]SiteLink{
			"svwikivoyage": SiteLink{
				Site:   "svwikivoyage",
				Title:  "Nordirland",
				Badges: []string{"test"},
			},
		},
		Claims: map[string][]Claim{
			claimCode: []Claim{
				Claim{
					ID:   "Q26$359b4cbf-41c5-4717-589c-eaee1ee2a323",
					Type: "statement",
					Rank: "normal",
					Mainsnak: Snak{
						Property: "P1549",
						Snaktype: "value",
						DataValue: DValue{
							Type: "monolingualtext",
							Value: MonolinqualTextDataValue{
								Text:     "Northern Irish",
								Language: "en",
							},
						},
					},
				},
			},
		},
	}
}
func TestTransform(t *testing.T) {
	logger := zerolog.New(os.Stdout)
	itemchan := make(chan WikiDataEntity, 1)
	transformer := MustNewTransform(logger, itemchan)
	tests := []struct {
		item           WikiDataEntity
		expectedOutput bool
	}{
		{buildWikiDataEntity("P10"), false},
		{buildWikiDataEntity("P214"), true},
	}
	for _, test := range tests {
		actual := transformer.Transform(test.item)
		assert.Equal(t, test.expectedOutput, actual)
	}
}
