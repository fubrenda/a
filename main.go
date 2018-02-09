package main

import (
	"log"
	"net/http"

	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt" // Makes sure bolt is included
	"github.com/voidfiles/a/api"
	"github.com/voidfiles/a/authority"
	"github.com/voidfiles/a/cli"
	"github.com/voidfiles/a/search"
)

func main() {
	args := cli.GetArgs()
	qs, err := graph.NewQuadStore(args.Db, args.Dbpath, graph.Options{"nosync": args.Nosync})
	if err != nil {
		panic(err)
	}
	index := search.MustNewIndex(args.IndexPath)

	resolver := authority.NewResolver(qs, index)

	log.Printf("Starting up an http server")
	mux := api.NewApi(resolver)

	address := args.IP + ":" + args.Port
	log.Fatal(http.ListenAndServe(address, mux))

}
