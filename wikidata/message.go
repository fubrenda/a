package wikidata

import (
	"encoding/json"
	"log"
)

type MessageStream struct {
	chunkSize      int
	data           *json.Decoder
	done           bool
	currentResults []WikiDataMessage
}

func (ms *MessageStream) Next() bool {
	if ms.done {
		return false
	}
	ms.currentResults = make([]WikiDataMessage, 0)

	for ms.data.More() {
		var msg WikiDataMessage
		// decode an array value (Message)
		err := ms.data.Decode(&msg)
		if err != nil {
			log.Fatal(err)
		}
		ms.currentResults = append(ms.currentResults, msg)
		if len(ms.currentResults) >= ms.chunkSize {
			return true
		}
	}

	ms.done = true

	return false
}

func (ms *MessageStream) Value() []WikiDataMessage {
	return ms.currentResults
}
