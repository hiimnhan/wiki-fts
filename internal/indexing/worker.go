package indexing

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hiimnhan/wiki-fts/common"
)

type Worker struct {
	records  common.Records
	online   bool
	listener chan *Msg
	sender   chan *Msg
	id       int
}

func NewWorker(listener chan *Msg, sender chan *Msg, id int) *Worker {
	return &Worker{
		records:  make(common.Records),
		listener: listener,
		sender:   sender,
		id:       id,
		online:   false,
	}
}

func (w *Worker) Run() {
	w.online = true
	for {
		select {
		case msg := <-w.listener:
			switch msg.Type {
			case MsgIndex:
				docs := msg.Data.(common.Documents)
				log.Infof("worker id %d starts indexing docs with len %d", w.id, len(docs))
				start := time.Now()
				w.index(docs)
				log.Infof("worker id %d finished index docs with time %v", w.id, time.Since(start))
			case MsgCombine:
				records := msg.Data.(common.Records)
				log.Infof("worker id %d starts combining docs with len %d", w.id, len(records))
				start := time.Now()
				log.Infof("worker id %d finished combining docs with time %v", w.id, time.Since(start))
				w.combine(records)
			case MsgRetireWorker:
				log.Info("Worker retire, id: ", w.id)
				w.online = false
			case MsgDeliverData:
				log.Infof("worker id %d starts delivering data", w.id)
				w.sender <- NewMsgDeliverData(w.records, w.id)
			}
		}
	}
}

func (w *Worker) index(docs common.Documents) {
	records := make(common.Records)
	for _, doc := range docs {
		for _, token := range common.TokenizeAndFilter(fmt.Sprintf("%s %s", doc.Text, doc.Title)) {
			set := records[token]
			if set == nil {
				set = &common.Set{}
			}
			set.Add(doc.ID)
		}
	}
}

func (w *Worker) combine(records common.Records) {
	for word, set := range records {
		for _, id := range set.Items() {
			if !w.records[word].Has(id) {
				w.records[word].Add(id)
			}
		}
	}
}
