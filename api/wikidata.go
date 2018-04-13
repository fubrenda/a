package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fubrenda/a/wikidata"
	"github.com/rs/zerolog"
)

type WikidataServer struct {
	logger zerolog.Logger
	db     *wikidata.WikiDataStore
}

func MustNewWikidataServer(logger zerolog.Logger, db *wikidata.WikiDataStore) *WikidataServer {

	ws := &WikidataServer{
		logger: logger,
		db:     db,
	}

	return ws
}

func httpError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	data := fmt.Sprintf("{\"error\": \"%s\"}", error)
	fmt.Fprintln(w, data)
}

func (ws *WikidataServer) FindByID(w http.ResponseWriter, r *http.Request) {
	b, err := ws.db.FindByIdentifier(r.URL.Query().Get("id"))
	if err != nil {
		httpError(w, "Server error", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(b)
	if err != nil {
		httpError(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (ws *WikidataServer) List(w http.ResponseWriter, r *http.Request) {
	b := ws.db.Scan([]byte(wikidata.IdentifierKeyPrefix), []byte(""), 200)

	data, err := json.Marshal(b)
	if err != nil {
		httpError(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
