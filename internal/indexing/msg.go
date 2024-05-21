package indexing

import "github.com/hiimnhan/wiki-fts/common"

type Msg struct {
	Data any
	Type int
	ID   int
}

func NewMsgIndex(docs common.Documents, id int) *Msg {
	return &Msg{Data: docs, Type: MsgIndex, ID: id}
}

func NewMsgCombine(records []common.Records) *Msg {
	return &Msg{Data: records, Type: MsgCombine}
}

func NewMsgRetireWorker() *Msg {
	return &Msg{Type: MsgRetireWorker}
}

func NewMsgDeliverData(records common.Records, id int) *Msg {
	return &Msg{Type: MsgDeliverData, Data: records, ID: id}
}

func NewMsgRequestData() *Msg {
	return &Msg{Type: MsgWorkerDelivery}
}

func NewMsgHealthcheck() *Msg {
	return &Msg{Type: MsgHealthcheck}
}

func NewMsgWorkerInfo(info Info, id int) *Msg {
	return &Msg{Type: MsgWorkerInfo, Data: info, ID: id}
}

func NewMsgWorkerCombineCompleted(id int) *Msg {
	return &Msg{Type: MsgWorkerCombineCompleted, ID: id}
}

func NewMsgWorkerCompleted(id int) *Msg {
	return &Msg{Type: MsgWorkerCompleted, ID: id}
}

func NewMsgSaveToDisk() *Msg {
	return &Msg{Type: MsgSaveToDisk}
}
