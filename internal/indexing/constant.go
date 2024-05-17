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
