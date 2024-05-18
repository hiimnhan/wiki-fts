package common

import (
	"compress/gzip"
	"encoding/xml"
	"os"

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

// load document
func LoadDocuments(path string) (Documents, error) {
	log.Infof("Loading documents %s...", path)
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
	}

	return docs, nil
}

type Record map[string]int
type Records map[string]*Set // word:ID of Document
