package wikidata

import "encoding/json"

type LangValue struct {
	Language string `json:"language"`
	Value    string `json:"value"`
}

type SiteLink struct {
	Badges []string `json:"badges"`
	Site   string   `json:"site"`
	Title  string   `json:"title"`
}

type TypeValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type MainSnak struct {
	DataType  string    `json:"datatype"`
	DataValue TypeValue `json:"datavalue"`
	Property  string    `json:"property"`
	Snaktype  string    `json:"snaktype"`
}

type Predicate struct {
	ID       string   `json:"id"`
	MainSnak MainSnak `json:"mainsnak"`
	Rank     string   `json:"rank"`
	Type     string   `json:"type"`
}

type Claim map[string]*json.RawMessage

type AliasMap map[string][]LangValue
type ClaimMap map[string]*json.RawMessage
type DescriptionMap map[string]LangValue

//WikiDataMessage is a generic interface to unmarshal unknown json to
type WikiDataMessage struct {
	Aliases      AliasMap             `json:"aliases"`
	Claims       ClaimMap             `json:"claims"`
	Descriptions DescriptionMap       `json:"descriptions"`
	Labels       map[string]LangValue `json:"lables"`
	SiteLinks    map[string]SiteLink  `json:"sitelinks"`
	Type         string               `json:"item"`
	ID           string               `json:"id"`
}
