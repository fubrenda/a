package wikidata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var JSONData = []byte(`{
	"id": "Q26",
	"type": "item",
	"labels": {
		"en-gb": {
			"language": "en-gb",
			"value": "Northern Ireland"
		}
	},
	"aliases": {
		"ca": [{
			"language": "ca",
			"value": "Ulster"
		}]
	},
	"claims": {
		"P1549": [{
				"mainsnak": {
					"snaktype": "value",
					"property": "P1549",
					"datavalue": {
						"value": {
							"text": "Northern Irish",
							"language": "en"
						},
						"type": "monolingualtext"
					},
					"datatype": "monolingualtext"
				},
				"type": "statement",
				"id": "Q26$359b4cbf-41c5-4717-589c-eaee1ee2a323",
				"rank": "normal"
			},
			{
				"mainsnak": {
					"snaktype": "value",
					"property": "P1549",
					"datavalue": {
						"value": {
							"text": "nordirisch",
							"language": "de"
						},
						"type": "monolingualtext"
					},
					"datatype": "monolingualtext"
				},
				"type": "statement",
				"id": "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				"rank": "normal"
			},
			{
				"mainsnak": {
					"snaktype": "value",
					"property": "P1549",
					"datavalue": {
						"value": "test",
						"type": "string"
					},
					"datatype": "string"
				},
				"type": "statement",
				"id": "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				"rank": "normal"
			},
			{
				"mainsnak": {
					"snaktype": "value",
					"property": "P1549",
					"datavalue": {
            "value": {
              "entity-type": "item",
              "numeric-id": 30
            },
						"type": "wikibase-entityid"
					},
					"datatype": "wikibase-entityid"
				},
				"type": "statement",
				"id": "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				"rank": "normal"
			},
			{
				"mainsnak": {
					"snaktype": "value",
					"property": "P1549",
          "datavalue": {
            "value": {
              "latitude": 52.516666666667,
              "longitude": 13.383333333333,
              "altitude": null,
              "precision": 0.016666666666667,
              "globe": "http:\/\/www.wikidata.org\/entity\/Q2"
            },
            "type": "globecoordinate"
          },
					"datatype": "globecoordinate"
				},
				"type": "statement",
				"id": "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				"rank": "normal"
			},
			{
				"mainsnak": {
					"snaktype": "value",
					"property": "P1549",
          "datavalue": {
            "value":{
              "amount":"+10.38",
              "upperBound":"+10.375",
              "lowerBound":"+10.385",
              "unit":"http://www.wikidata.org/entity/Q712226"
            },
            "type":"quantity"
          },
					"datatype": "quantity"
				},
				"type": "statement",
				"id": "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				"rank": "normal"
			},
			{
				"mainsnak": {
					"snaktype": "value",
					"property": "P1549",
          "datavalue": {
            "value": {
              "time": "+2001-12-31T00:00:00Z",
              "timezone": 0,
              "before": 0,
              "after": 0,
              "precision": 11,
              "calendarmodel": "http:\/\/www.wikidata.org\/entity\/Q1985727"
            },
            "type": "time"
          },
					"datatype": "time"
				},
				"type": "statement",
				"id": "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				"rank": "normal"
			}
		]
	},
	"sitelinks": {
		"svwikivoyage": {
			"site": "svwikivoyage",
			"title": "Nordirland",
			"badges": ["test"]
		}
	}
}`)

var CaUlster = LanguageValue{
	Language: "ca",
	Value:    "Ulster",
}

var ValidWikiDataEntity = WikiDataEntity{
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
		"P1549": []Claim{
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
			Claim{
				ID:   "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				Type: "statement",
				Rank: "normal",
				Mainsnak: Snak{
					Property: "P1549",
					Snaktype: "value",
					DataValue: DValue{
						Type: "monolingualtext",
						Value: MonolinqualTextDataValue{
							Text:     "nordirisch",
							Language: "de",
						},
					},
				},
			},
			Claim{
				ID:   "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				Type: "statement",
				Rank: "normal",
				Mainsnak: Snak{
					Property: "P1549",
					Snaktype: "value",
					DataValue: DValue{
						Type:  "string",
						Value: "test",
					},
				},
			},
			Claim{
				ID:   "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				Type: "statement",
				Rank: "normal",
				Mainsnak: Snak{
					Property: "P1549",
					Snaktype: "value",
					DataValue: DValue{
						Type: "wikibase-entityid",
						Value: WikiBaseEntityIDDataValue{
							EntityType: "item",
							NumericID:  uint64(30),
						},
					},
				},
			},
			Claim{
				ID:   "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				Type: "statement",
				Rank: "normal",
				Mainsnak: Snak{
					Property: "P1549",
					Snaktype: "value",
					DataValue: DValue{
						Type: "globecoordinate",
						Value: GlobeCoordinateDataValue{
							Latitude:  float32(52.516666666667),
							Longitude: float32(13.383333333333),
							Pecision:  float32(0.016666666666667),
							Globe:     "http://www.wikidata.org/entity/Q2",
						},
					},
				},
			},
			Claim{
				ID:   "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				Type: "statement",
				Rank: "normal",
				Mainsnak: Snak{
					Property: "P1549",
					Snaktype: "value",
					DataValue: DValue{
						Type: "quantity",
						Value: QuantityDataValue{
							Amount:     "+10.38",
							UpperBound: "+10.375",
							LowerBound: "+10.385",
							Unit:       "http://www.wikidata.org/entity/Q712226",
						},
					},
				},
			},
			Claim{
				ID:   "Q26$516ef7a7-48f8-5bdc-418e-d11f149b3759",
				Type: "statement",
				Rank: "normal",
				Mainsnak: Snak{
					Property: "P1549",
					Snaktype: "value",
					DataValue: DValue{
						Type: "time",
						Value: TimeDataValue{
							Time:          "+2001-12-31T00:00:00Z",
							Timezone:      0,
							Before:        int64(0),
							After:         int64(0),
							Percision:     0,
							CalenderModel: "http://www.wikidata.org/entity/Q1985727",
						},
					},
				},
			},
		},
	},
}

func TestWikiDataEntity(t *testing.T) {
	var Entity = &WikiDataEntity{}
	json.Unmarshal(JSONData, Entity)
	assert.Equal(t, ValidWikiDataEntity, *Entity)
}
