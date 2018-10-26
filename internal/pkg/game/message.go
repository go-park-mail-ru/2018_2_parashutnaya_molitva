package game

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

const (
	TurnMsg  = "turn"
	ErrorMsg = "error"
)

var (
	errDataIsEmpty              = errors.New("Data is empty")
	errImpossibleUnmarshalToMsg = errors.New("Impossible unmarshal to Message")
)

type Message struct {
	msgType string
	data    json.RawMessage
}

func UnmarshalToMessage(data []byte) (*Message, error) {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, errImpossibleUnmarshalToMsg
	}

	return msg, nil
}

func (m *Message) ToUnmarshalData(v interface{}) error {
	// Не возвращает ошибку
	data, _ := m.data.MarshalJSON()
	if reflect.DeepEqual(data, []byte("null")) {
		return errDataIsEmpty
	}

	return json.Unmarshal(data, v)
}

func NewMessage(msgType string, data []byte) *Message {
	return &Message{
		msgType: msgType,
		data:    json.RawMessage(data),
	}
}
