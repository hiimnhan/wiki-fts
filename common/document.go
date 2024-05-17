package common

import (
	"compress/gzip"
	"encoding/xml"
	"os"
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

// load document
func LoadDocuments(path string) ([]Document, error) {
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
