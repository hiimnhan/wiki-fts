package common

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type Sublink struct {
	Anchor string `xml:"anchor"`
	URL    string `xml:"link"`
}

type Links struct {
	Sublinks []Sublink `xml:"sublink"`
}

type Document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	Links Links  `xml:"links"`
	ID    int
}

type Documents []Document
type DocumentDict map[int]Document

// load document
func LoadDocuments(path string) (Documents, error) {
	log.Infof("Loading documents %s...", path)
	start := time.Now()
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, NewError("Documents", err)
	}
	defer gz.Close()

	decoder := xml.NewDecoder(gz)

	d := struct {
		Documents []Document `xml:"doc"`
	}{}

	if err := decoder.Decode(&d); err != nil {
		return nil, NewError("Documents", err)
	}

	docs := d.Documents
	for i := range docs {
		docs[i].ID = i
		after, found := strings.CutPrefix(docs[i].Title, "wikipedia: ")
		if found {
			docs[i].Title = after
		}
	}

	log.Infof("Loaded %d documents in %v", len(docs), time.Since(start))
	return docs, nil
}

func (d *Documents) SaveDocsDictToDisk(path string) error {
	docDict := make(DocumentDict)
	log.Infof("Saving documents to %s...", path)
	for _, doc := range *d {
		docDict[doc.ID] = doc
	}
	b, err := json.Marshal(docDict)
	if err != nil {
		return NewError("documents", err)
	}

	err = os.WriteFile(path, b, 0644)
	if err != nil {
		return NewError("documents", err)
	}
	log.Infof("Saved %d documents to %s", len(docDict), path)
	return nil
}

type Record map[string]int
type Records map[string]*Set // word:ID of Document
