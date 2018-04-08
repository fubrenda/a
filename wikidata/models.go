package wikidata

import (
	"encoding/json"
)

//MonolinqualTextDataValue is a field for a string that is not translated into other languages.
type MonolinqualTextDataValue struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

//WikiBaseEntityIDDataValue are used to reference entities on the same repository.
type WikiBaseEntityIDDataValue struct {
	EntityType string `json:"entity-type"`
	NumericID  uint64 `json:"numeric-id"`
}

type GlobeCoordinateDataValue struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Pecision  float32 `json:"precision"`
	Globe     string  `json:"globe"`
}

//QuantityDataValue is a quantity that relates to some kind of well-defined unit.
type QuantityDataValue struct {
	Amount     string `json:"amount"`
	UpperBound string `json:"upperBound"`
	LowerBound string `json:"lowerBound"`
	Unit       string `json:"unit"`
}

//TimeDataValue is a Literal data field for a point in time
type TimeDataValue struct {
	Time          string `json:"time"`
	Timezone      int    `json:"timezone"`
	Before        int64  `json:"before"`
	After         int64  `json:"after"`
	Percision     int    `json:"percision"`
	CalenderModel string `json:"calendarmodel"`
}

type innerDValue struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

//DValue is a field that contains the actual value the Snak associates with the Property
type DValue struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (d *DValue) UnmarshalJSON(data []byte) error {
	var inner = innerDValue{}
	json.Unmarshal(data, &inner)

	d.Type = inner.Type
	switch d.Type {
	case "monolingualtext":
		val := MonolinqualTextDataValue{}
		json.Unmarshal(inner.Value, &val)
		d.Value = val
	case "string":
		var val string
		json.Unmarshal(inner.Value, &val)
		d.Value = val
	case "wikibase-entityid":
		val := WikiBaseEntityIDDataValue{}
		json.Unmarshal(inner.Value, &val)
		d.Value = val
	case "globecoordinate":
		val := GlobeCoordinateDataValue{}
		json.Unmarshal(inner.Value, &val)
		d.Value = val
	case "time":
		val := TimeDataValue{}
		json.Unmarshal(inner.Value, &val)
		d.Value = val
	case "quantity":
		val := QuantityDataValue{}
		json.Unmarshal(inner.Value, &val)
		d.Value = val
	}

	return nil
}

//Snak provides some kind of information about a specific Property of a given Entity
type Snak struct {
	Property  string `json:"property"`
	Snaktype  string `json:"snaktype"`
	DataValue DValue `json:"datavalue"`
}

type QualifierSnak struct {
	Snak
	Hash string `json:"hash"`
}

type ReferenceSnaks struct {
	Hash  string             `json:"hash"`
	Snaks *map[string][]Snak `json:"snaks"`
}

type References []ReferenceSnaks

//Claim consists of a main value (or main Snak) and a number of qualifier Snaks
type Claim struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Rank       string `json:"rank"`
	Mainsnak   Snak   `json:"mainsnak"`
	Qualifiers *map[string][]QualifierSnak
	References *References
}

//LanguageValue represents the language code and value of a label or descripition
type LanguageValue struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

//SiteLink are given as records for each site identifier.
type SiteLink struct {
	Site   string   `json:"site"`
	Title  string   `json:"title"`
	Badges []string `json:"badges"`
}

//LabelMap is a set of labels
type LabelMap map[string]LanguageValue

//LabelMap is a set of aliases
type AliasMap map[string][]LanguageValue

//SiteLinkMap is a set of site links
type SiteLinkMap map[string]SiteLink

//ClaimMap is a set of Claims
type ClaimMap map[string][]Claim

//WikiDataEntity is a struct that represents a WikiDataEntity from a Wikidata dump
type WikiDataEntity struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Labels    LabelMap    `json:"labels"`
	Aliases   AliasMap    `json:"aliases"`
	SiteLinks SiteLinkMap `json:"sitelinks"`
	Claims    ClaimMap    `json:"claims"`
}
