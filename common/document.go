package common

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"syscall"
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

func (d *Document) Display() string {
	var s string
	s += fmt.Sprintf("Title: %s\n", d.Title)
	s += fmt.Sprintf("URL: %s\n", d.URL)
	s += fmt.Sprintf("Text: %s\n\n", d.Text)
	s += "Links:\n"
	for _, link := range d.Links.Sublinks {
		s += fmt.Sprintf("%s: %s\n", link.Anchor, link.URL)
	}

	return s
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

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	log.Printf("File size %d", size)
	if size <= 0 || size != int64(int(size)) {
		return nil, fmt.Errorf("Invalid file size %d", size)
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	defer syscall.Munmap(data)

	d := struct {
		Documents []Document `xml:"doc"`
	}{}

	err = xml.Unmarshal(data, &d)
	if err != nil {
		return nil, err
	}

	docs := d.Documents
	for i := range docs {
		docs[i].ID = i
		after, found := strings.CutPrefix(docs[i].Title, "Wikipedia: ")
		if found {
			docs[i].Title = after
		}
	}

	log.Infof("Loaded %d documents in %v", len(docs), time.Since(start))
	return docs, nil
}

func (d *Documents) GenerateDocsDictionary() *DocumentDict {
	docDict := make(DocumentDict)
	for _, doc := range *d {
		docDict[doc.ID] = doc
	}
	return &docDict
}

type Record map[string]int
type Records map[string]*Set // word:ID of Document

type Index map[string][]int

func (idx *Index) FindIndexes(text string) []int {
	log.Printf("Querying %q\n....", text)
	var res [][]int
	for _, word := range TokenizeAndFilter(text) {
		log.Printf("Word %q...", word)
		if ids, ok := (*idx)[word]; ok {
			res = append(res, ids)
		} else {
			res = append(res, []int{})
		}
	}

	return Intersect(res)
}
