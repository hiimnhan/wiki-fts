package indexing

type Master struct {
	workers      map[int]chan *Msg
	listener     chan *Msg
	numOfWorkers int
	id           int
}

func NewMaster(numOfWorkers int) *Master {
	return &Master{
		workers:      make(map[int]chan *Msg),
		listener:     make(chan *Msg),
		numOfWorkers: numOfWorkers,
		id:           0,
	}
}

func (m *Master) newWorker(id int)
