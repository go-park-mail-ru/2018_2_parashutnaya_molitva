package game

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

const (
	TurnMsg = "turn"
)

var (
	errDataIsEmpty = errors.New("Data is empty")
)

type Message struct {
	msgType string
	data    json.RawMessage
}

func (m *Message) ToUnmarshalData(v interface{}) error {
	// Не возвращает ошибку
	data, _ := m.data.MarshalJSON()
	if reflect.DeepEqual(data, []byte("null")) {
		return errDataIsEmpty
	}

	return json.Unmarshal(data, v)
}
