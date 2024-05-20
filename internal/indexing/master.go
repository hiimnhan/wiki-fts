package indexing

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hiimnhan/wiki-fts/common"
)

type Master struct {
	workers      map[int]chan *Msg
	online       map[int]bool
	available    []int
	listener     chan *Msg
	numOfWorkers int
	id           int
	mu           sync.RWMutex
	running      bool
	records      []common.Records
}

func NewMaster(numOfWorkers int) *Master {
	return &Master{
		workers:      make(map[int]chan *Msg),
		online:       make(map[int]bool),
		available:    make([]int, 0, numOfWorkers),
		listener:     make(chan *Msg, MasterChanCap),
		numOfWorkers: numOfWorkers,
		id:           0,
		running:      true,
		records:      make([]common.Records, 0),
	}
}

func (m *Master) Run(path string) {
	for id := 1; id <= m.numOfWorkers; id++ {
		worker, err := m.newWorker(id)
		if err != nil {
			log.Error(err)
			continue
		}
		go worker.Run()
	}

	numsOfMapWorkers := m.numOfWorkers / 2
	workloads, err := m.generateWorkloads(path, numsOfMapWorkers)
	if err != nil {
		log.Fatalf("Cannot generate workload %v", err)
	}

	err = m.delegateInitialWorkload(workloads)
	if err != nil {
		log.Fatalf("Cannot delegate workload %v", err)
	}

	reduceCompleted := 0

	for m.running {
		select {
		case msg := <-m.listener:
			switch msg.Type {
			case MsgWorkerCompleted:
				workerId := msg.ID
				records := msg.Data.(common.Records)
				m.records = append(m.records, records)
				log.Infof("worker %d finished tasks, len %d", workerId, len(records))
				m.retireWorker(workerId)
				reduceCompleted++
			case MsgDeliverData:
				records := msg.Data.(common.Records)
				workerId := msg.ID
				log.Infof("master receives data from worker %d, len %d", workerId, len(records))
				m.transferData(workerId, records)
			}
		case <-time.After(2 * time.Second):
			log.Info("master idle, checking status workers...")

			onlineWorkers, size := m.onlineWorkers()
			switch size {
			case 0:
				log.Warn("no worker online")
			default:
				if reduceCompleted == len(workloads) {
					for _, w := range onlineWorkers {
						m.retireWorker(w)
					}
					m.saveRecordsToDisk()
					m.running = false
				}
			}
		}
	}

	log.Warn("master stopped")
	return
}

func (m *Master) newWorker(id int) (*Worker, error) {
	if _, exist := m.workers[id]; exist {
		return nil, common.NewError("master", fmt.Errorf("worker with id %d already existed", id))
	}
	channel := make(chan *Msg, WorkerChanCap)
	worker := NewWorker(channel, m.listener, id)
	m.workers[id] = channel
	m.online[id] = false
	m.available = append(m.available, id)

	return worker, nil
}

func (m *Master) retireWorker(id int) {
	m.workers[id] <- NewMsgRetireWorker()
	m.online[id] = false
}

func (m *Master) generateWorkloads(path string, numOfWorkers int) ([]common.Documents, error) {
	docs, err := common.LoadDocuments(path)
	if err != nil {
		return nil, err
	}
	if numOfWorkers == 0 {
		numOfWorkers = DEFAULT_NUM_WORKERS
	}

	var workloads []common.Documents
	chunkSize := len(docs) / numOfWorkers
	for i := 0; i < len(docs); i += chunkSize {
		end := i + chunkSize
		if end > len(docs) {
			end = len(docs)
		}
		workloads = append(workloads, docs[i:end])
	}

	return workloads, nil
}

func (m *Master) nextAvailableWorker() (int, error) {
	if len(m.available) == 0 {
		return 0, common.NewError("master", errors.New("no worker available"))
	}
	log.Warnf("available %v", m.available)
	m.mu.Lock()
	next, _ := common.Shift(&m.available)
	m.mu.Unlock()
	m.online[next] = true
	log.Warnf("next worker %d", next)

	return next, nil
}

func (m *Master) delegateInitialWorkload(workloads []common.Documents) error {
	for _, workload := range workloads {
		workerId, err := m.nextAvailableWorker()
		if err != nil {
			return err
		}
		m.workers[workerId] <- NewMsgIndex(workload, m.id)
	}

	return nil
}

func (m *Master) saveRecordsToDisk() {
	// save records to disk
	b, err := json.Marshal(m.records)
	if err != nil {
		log.Fatalf("Cannot marshal records %v", err)
	}

	err = os.WriteFile(common.OutputPath, b, 0644)
	if err != nil {
		log.Fatalf("Cannot write records to disk %v", err)
	}

}

func (m *Master) transferData(prevWorker int, data common.Records) {
	worker, err := m.nextAvailableWorker()
	if err != nil {
		log.Fatalf("Can not move to next stage %v", err)
	}

	log.Infof("redirect data from worker %d to worker %d", prevWorker, worker)
	m.workers[worker] <- NewMsgCombine(data, prevWorker)
	m.available = append(m.available, prevWorker)
}

func (m *Master) onlineWorkers() (workers []int, size int) {
	for i := range m.online {
		if m.online[i] {
			workers = append(workers, i)
			size++
		}
	}

	return
}

func (m *Master) healthcheck() {
	workers, _ := m.onlineWorkers()
	for _, w := range workers {
		log.Infof("healthcheck worker %d...", w)
		m.workers[w] <- NewMsgHealthcheck()
	}
}
