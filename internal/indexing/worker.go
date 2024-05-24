package indexing

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hiimnhan/wiki-fts/common"
)

type Stage string

const (
	OFFLINE Stage = "offline"
	MAP     Stage = "map"
	REDUCE  Stage = "reduce"
)

type State int

const (
	IDLE        State = 0
	IN_PROGRESS State = 1
	COMPLETED   State = 2
)

type Info struct {
	workState State
	stage     Stage
}

type Worker struct {
	records  common.Records
	online   bool
	listener chan *Msg
	sender   chan *Msg
	id       int
	info     Info
}

func NewWorker(listener chan *Msg, sender chan *Msg, id int) *Worker {
	return &Worker{
		records:  make(common.Records),
		listener: listener,
		sender:   sender,
		id:       id,
		online:   false,
		info:     Info{workState: IDLE, stage: OFFLINE},
	}
}

func (w *Worker) Run() {
	w.online = true
	for w.online {
		select {
		case msg := <-w.listener:
			switch msg.Type {
			case MsgIndex:
				w.info.stage = MAP
				docs := msg.Data.(common.Documents)
				log.Infof("worker id %d starts indexing docs with len %d", w.id, len(docs))
				w.info.workState = IN_PROGRESS
				start := time.Now()
				w.index(docs)
				w.info.workState = COMPLETED
				log.Infof("worker id %d finished index docs with time %v, len %d", w.id, time.Since(start), len(w.records))
				log.Infof("worker id %d starts delivering data", w.id)
				w.sender <- NewMsgDeliverData(w.records, w.id)
				w.records = make(common.Records)
			case MsgCombine:
				w.info.stage = REDUCE
				records := msg.Data.([]common.Records)
				w.info.workState = IN_PROGRESS
				log.Infof("worker %d received data from master, len %d", w.id, len(records))
				log.Infof("worker id %d starts combining docs with len %d", w.id, len(records))
				start := time.Now()
				w.combine(records)
				w.info.workState = COMPLETED
				log.Infof("worker id %d finished combining docs with time %v", w.id, time.Since(start))
				log.Infof("worker id %d starts delivering data", w.id)
				w.sender <- NewMsgWorkerCombineCompleted(w.id)
			case MsgRetireWorker:
				w.info.stage = OFFLINE
				log.Infof("Worker retire, id: %d", w.id)
				w.online = false
				w.info.workState = IDLE
			case MsgSaveToDisk:
				log.Infof("worker %d starts saving records to disk", w.id)
				w.saveRecordsToDisk()
				log.Infof("worker %d finished saving records to disk", w.id)
				w.sender <- NewMsgWorkerCompleted(w.id)
			case MsgHealthcheck:
				log.Infof("worker %d send info to master", w.id)
				w.sender <- NewMsgWorkerInfo(w.info, w.id)
			}
		case <-time.After(2 * time.Second):
			continue
		}
	}

	log.Warnf("worker %d stopped", w.id)
}

func (w *Worker) index(docs common.Documents) {
	log.Warnf("worker %d indexing...", w.id)

	records := make(common.Records)
	for _, doc := range docs {
		for _, token := range common.TokenizeAndFilter(fmt.Sprintf("%s %s", doc.Text, doc.Title)) {
			set := records[token]
			if set == nil {
				set = new(common.Set)
				records[token] = set
			}
			set.Add(doc.ID)
		}
	}

	w.records = records

	log.Infof("worker %d finished indexing...", w.id)

}

func (w *Worker) combine(records []common.Records) {
	// combine records from other workers, return a map of token -> set of doc ids
	combinedRecords := make(common.Records)
	for _, r := range records {
		for token, set := range r {
			if combinedSet, exist := combinedRecords[token]; exist {
				combinedSet.Merge(set)
			} else {
				combinedRecords[token] = set
			}
		}
	}
	w.records = combinedRecords
	log.Infof("worker %d finished combining..., len %d", w.id, len(combinedRecords))
}

func (w *Worker) saveRecordsToDisk() {
	// save records to disk
	invertedIndex := make(common.Index)
	for token, set := range w.records {
		invertedIndex[token] = set.Items()
	}
	b, err := json.Marshal(invertedIndex)
	if err != nil {
		log.Fatalf("Cannot marshal records %v", err)
	}

	err = os.WriteFile(common.OutputPath, b, 0644)
	if err != nil {
		log.Fatalf("Cannot write records to disk %v", err)
	}
	log.Infof("worker %d saved records to disk", w.id)
}
