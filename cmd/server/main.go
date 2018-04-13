package main

import (
	"net/http"
	"os"
	"time"

	"github.com/coreos/bbolt"
	"github.com/fubrenda/a/api"
	"github.com/fubrenda/a/cli"
	"github.com/fubrenda/a/wikidata"
	"github.com/rs/zerolog"
)

type RequestLogger struct {
	logger  zerolog.Logger
	handler http.Handler
}

func NewRequestLogger(logger zerolog.Logger, handler http.Handler) *RequestLogger {
	return &RequestLogger{
		logger:  logger,
		handler: handler,
	}
}

func (rl *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rl.logger.Info().
		Str("scheme", r.URL.Scheme).
		Str("host", r.URL.Host).
		Str("path", r.URL.Path).
		Msg("logged request")
	rl.handler.ServeHTTP(w, r)
}

func main() {
	args := cli.GetArgs()

	wikidatadb, err := bolt.Open(args.WikidataDB, 0666, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})
	defer wikidatadb.Close()
	if err != nil {
		panic(err)
	}

	wikidatadbStore := wikidata.MustNewWikiDataStore(wikidatadb)

	db, err := bolt.Open(args.Dbpath, 0666, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})
	defer db.Close()

	if err != nil {
		panic(err)
	}
	logger := zerolog.New(os.Stdout)
	handler := http.NewServeMux()
	//recordStore := recordstore.MustNewRecordStore(db)
	wds := api.MustNewWikidataServer(logger, wikidatadbStore)
	handler.HandleFunc("/wikidata/list/", wds.List)
	handler.HandleFunc("/wikidata/by/id/", wds.FindByID)
	loggedHandler := NewRequestLogger(logger, handler)

	http.ListenAndServe("0.0.0.0:8081", loggedHandler)
}
