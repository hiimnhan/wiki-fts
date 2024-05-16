package common

import "fmt"

type ErrMsg struct {
	msg       string
	component string
}

func (e *ErrMsg) Error() string {
	return fmt.Sprintf("%s::%s", e.component, e.msg)

}

func NewError(component string, err error) *ErrMsg {
	return &ErrMsg{component, err.Error()}
}
