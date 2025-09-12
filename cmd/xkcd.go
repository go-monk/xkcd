package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/go-monk/xkcd"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("xkcd: ")

	maxConcurrency := flag.Int("c", 20, "max number of concurrent http requests when building offline index")
	indexFile := flag.String("i", "xkcd.json", "file holding offline index of comics")
	printTranscript := flag.Bool("t", false, "print also the transcript")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("supply a search term")
	}
	term := strings.Join(flag.Args(), " ")

	index := xkcd.NewIndex(*indexFile)

	if _, err := os.Stat(*indexFile); os.IsNotExist(err) {
		log.Printf("building index at %q ...\n", *indexFile)
		err := index.Build(*maxConcurrency)
		if err != nil {
			log.Fatalf("bulding index: %v", err)
		}
	}

	comics, err := index.Search(term)
	if err != nil {
		log.Fatalf("searching index: %v", err)
	}

	sort.Slice(comics, func(i, j int) bool {
		return comics[i].Num < comics[j].Num
	})

	for _, c := range comics {
		url := fmt.Sprintf("https://xkcd.com/%d/", c.Num)
		fmt.Printf("#%4d (%s) %-35s %s\n", c.Num, c.Year, c.Title, url)
		if *printTranscript {
			fmt.Printf("%s\n", c.Transcript)
			fmt.Println("---")
		}
	}
}
