package indexing

const (
	WorkerChanCap = 512
	MasterChanCap = 512
)

const (
	MsgIndex = iota + 1
	MsgCombine
	MsgRetireWorker
	MsgDeliverData
	MsgClearData
	MsgSaveToDisk
	MsgWorkerCombineCompleted
	MsgWorkerCompleted
	MsgWorkerDelivery
	MsgHealthcheck
	MsgWorkerInfo
)

const (
	DEFAULT_NUM_WORKERS = 10
)
