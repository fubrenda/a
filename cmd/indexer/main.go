package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/boutros/marc"
	"github.com/coreos/bbolt"
	"github.com/fubrenda/a/cli"
	"github.com/fubrenda/a/lcsh"
	"github.com/fubrenda/a/pipeline"
	"github.com/fubrenda/a/recordstore"
)

func detectFormat(f *os.File) (marc.Format, error) {
	sniff := make([]byte, 64)
	_, err := f.Read(sniff)
	if err != nil {
		log.Fatal(err)
	}
	format := marc.DetectFormat(sniff)

	// rewind reader
	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	switch format {
	case marc.MARC, marc.LineMARC, marc.MARCXML:
		return format, nil
	default:
		return format, errors.New("unknown MARC format")
	}
}

func marcReader(inputPath string) (io.Reader, marc.Format) {

	if inputPath == "" {
		reader := bufio.NewReader(os.Stdin)
		format := marc.MARCXML

		return reader, format
	}

	reader, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}

	format, err := detectFormat(reader)

	if err != nil {
		panic(err)
	}

	return reader, format
}

func main() {
	args := cli.GetArgs()

	db, err := bolt.Open(args.Dbpath, 0666, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})
	defer db.Close()

	if err != nil {
		panic(err)
	}
	recordStore := recordstore.MustNewRecordStore(db)
	reader, format := marcReader(args.InputPath)

	pl := lcsh.NewFileToRecordStorePipeline(recordStore, reader, format)
	pipeline.RunReporter(pl)
	if stats, err := recordStore.Stats(); err != nil {
		log.Print(err)
	} else {
		log.Print(stats)
	}

}
