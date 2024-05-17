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

func NewMsgCombine(records common.Records, id int) *Msg {
	return &Msg{Data: records, Type: MsgCombine, ID: id}
}

func NewMsgRetireWorker() *Msg {
	return &Msg{Type: MsgRetireWorker}
}

func NewMsgDeliverData(records common.Records, id int) *Msg {
	return &Msg{Type: MsgDeliverData, Data: records, ID: id}
}
