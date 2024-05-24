package main

import (
	"encoding/json"
	"flag"
	"slices"

	// "fmt"
	"io"
	"log"
	"os"

	// "runtime"
	"time"

	"github.com/hiimnhan/wiki-fts/common"
	"github.com/hiimnhan/wiki-fts/internal/indexing"
)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()
	master := indexing.NewMaster(10)
	master.Run(common.WikiDumpPath)
	log.Printf("All processes took %v", time.Since(start))

	jsonFile, err := os.Open(common.OutputPath)
	if err != nil {
		log.Fatalf("Cannot open file %s", common.OutputPath)
	}
	defer jsonFile.Close()

	bytesVal, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Cannot read file %s", common.OutputPath)
	}

	var records = make(common.Index)

	json.Unmarshal(bytesVal, &records)

	var query string
	flag.StringVar(&query, "q", "", "search query")
	flag.Parse()

	start = time.Now()
	matchedIds := records.FindIndexes(query)

	slices.Sort(matchedIds)

	// for _, id := range matchedIds {
	// 	doc := (*docsDict)[id]
	// 	fmt.Println(doc.Display())
	// }

	log.Printf("Search took %v, found %d, %v", time.Since(start), len(matchedIds), matchedIds)
}
