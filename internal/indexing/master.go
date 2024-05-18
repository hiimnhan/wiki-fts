package indexing

import (
	"errors"
	"fmt"
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
}

func NewMaster(numOfWorkers int) *Master {
	return &Master{
		workers:      make(map[int]chan *Msg),
		online:       make(map[int]bool),
		available:    make([]int, 0, numOfWorkers),
		listener:     make(chan *Msg, MasterChanCap),
		numOfWorkers: numOfWorkers,
		id:           0,
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

	log.Warnf("num available %d, %v", len(m.available), m.available)

	numsOfMapWorkers := m.numOfWorkers / 2
	workloads, err := m.generateWorkloads(path, numsOfMapWorkers)
	if err != nil {
		log.Fatalf("Cannot generate workload %v", err)
	}

	err = m.delegateInitialWorkload(workloads)
	if err != nil {
		log.Fatalf("Cannot delegate workload %v", err)
	}

	deliverReceived := 0

	for {
		select {
		case msg := <-m.listener:
			switch msg.Type {
			case MsgWorkerCompleted:
				workerId := msg.ID
				log.Infof("worker id %d finished tasks", workerId)
				m.retireWorker(workerId)
				m.available = append(m.available, workerId)
			case MsgDeliverData:
				records := msg.Data.(common.Records)
				workerId := msg.ID
				log.Infof("master receives data from worker id %d", workerId)
				m.transferData(workerId, records)
				m.retireWorker(workerId)
				deliverReceived++
			}
		case <-time.After(2 * time.Second):
			log.Info("master idle, checking status workers...")

			var onlineWorkers int
			for i := 0; i < len(m.online); i++ {
				if m.online[i] {
					onlineWorkers++
				}
			}
			log.Infof("online workers %d", onlineWorkers)
			log.Infof("available workers %d", len(m.available))
		}
	}

}

func (m *Master) newWorker(id int) (*Worker, error) {
	if _, exist := m.workers[id]; exist {
		return nil, common.NewError("master", fmt.Errorf("worker with id %d already existed", id))
	}
	channel := make(chan *Msg, WorkerChanCap)
	worker := NewWorker(channel, m.listener, id)
	m.workers[id] = channel
	m.online[id] = true
	m.available = append(m.available, id)

	return worker, nil
}

func (m *Master) retireWorker(id int) {
	m.workers[id] <- NewMsgRetireWorker()
	m.online[id] = false
	m.available = append(m.available, id)
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
	next := common.Pop(&m.available)
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

func (m *Master) transferData(prevWorker int, data common.Records) {
	worker, err := m.nextAvailableWorker()
	if err != nil {
		log.Fatalf("Can not move to next stage %v", err)
	}

	log.Infof("redirect data from worker %d to worker %d", prevWorker, worker)
	m.workers[worker] <- NewMsgCombine(data, prevWorker)
}

// func (m *Master) requestData() {
// 	for id := range m.inused {
// 		if m.inused[id] {
// 			log.Infof("requesting data from worker %d...", id)
// 			m.workers[id] <- NewMsgRequestData()
// 		}
// 	}
// }
