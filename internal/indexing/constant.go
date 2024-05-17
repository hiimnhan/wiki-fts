package indexing

const (
	MsgIndex = iota + 1
	MsgCombine
	MsgRetireWorker
	MsgDeliverData
	MsgClearData
	MsgSortSave
	MsgWorkerCompleted
	MsgWorkerDelivery
)

const (
	DEFAULT_NUM_WORKERS = 10
)
