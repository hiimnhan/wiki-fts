package indexing

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/hiimnhan/wiki-fts/common"
)

type Master struct {
	workers      map[int]chan *Msg
	online       map[int]bool
	inused       map[int]bool
	listener     chan *Msg
	numOfWorkers int
	id           int
	workloads    []common.Documents
}

func NewMaster(numOfWorkers int) *Master {
	return &Master{
		workers:      make(map[int]chan *Msg),
		online:       make(map[int]bool),
		inused:       make(map[int]bool),
		listener:     make(chan *Msg),
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

	numsOfMapWorkers := m.numOfWorkers / 2
	m.generateWorkloads(path, numsOfMapWorkers)
	err := m.delegateInitialWorkload()
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
				m.inused[workerId] = false
			case MsgDeliverData:
				records := msg.Data.(common.Records)
				workerId := msg.ID
				log.Infof("master receives data from worker id %d", workerId)
				m.retireWorker(workerId)
				deliverReceived++
			}
		default:
			log.Info("master idle, checking status workers...")
			var onlineWorkers int
			for i := 0; i < len(m.online); i++ {
				if m.online[i] {
					onlineWorkers++
				}
			}
			log.Infof("online workers %d", onlineWorkers)

			var inusedWorkers int
			for i := 0; i < len(m.inused); i++ {
				if m.online[i] {
					inusedWorkers++
				}
			}
			log.Infof("online workers %d", inusedWorkers)

		}
	}

}

func (m *Master) newWorker(id int) (*Worker, error) {
	if _, exist := m.workers[id]; !exist {
		return nil, common.NewError("master", errors.New(fmt.Sprintf("worker with id %d already existed", id)))
	}
	channel := make(chan *Msg)
	worker := NewWorker(channel, m.listener, id)
	m.workers[id] = channel
	m.online[id] = true
	m.inused[id] = false

	return worker, nil
}

func (m *Master) retireWorker(id int) {
	m.workers[id] <- NewMsgRetireWorker()
	m.online[id] = false
	m.inused[id] = false
}

func (m *Master) generateWorkloads(path string, numOfWorkers int) error {
	docs, err := common.LoadDocuments(path)
	if err != nil {
		return err
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

	m.workloads = workloads
	return nil
}

func (m *Master) nextAvailableWorker() (int, error) {
	for inusedW := range m.inused {
		if m.inused[inusedW] == false {
			if m.online[inusedW] {
				return inusedW, nil
			}
		}
	}
	return 0, common.NewError("master", errors.New("no worker available"))
}

func (m *Master) delegateInitialWorkload() error {
	for _, workload := range m.workloads {
		workerId, err := m.nextAvailableWorker()
		if err != nil {
			return err
		}

		m.workers[workerId] <- NewMsgIndex(workload, m.id)
		m.inused[workerId] = true
	}

	return nil
}
