package game

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

const (
	TurnMsg   = "turn"
	ErrorMsg  = "error"
	InitMsg   = "init"
	InfoMsg   = "info"
	StartMsg  = "start"
	ResultMsg = "result"
)

var (
	errDataIsEmpty              = errors.New("Data is empty")
	errImpossibleUnmarshalToMsg = errors.New("Impossible unmarshal to Message")
)

type ResultMessage struct {
	Result string `jsong:"result"`
	Score  int    `json:"score"`
}

type TurnMessage struct {
	Turn string `json:"turn"`
}

type StartMessage struct {
	Color string `json:"color"`
}

type InitMessage struct {
	RoomId string `json:"roomid"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type InfoMessage struct {
	Info string `json:"info"`
}

type Message struct {
	MsgType string
	Data    json.RawMessage
}

func UnmarshalToMessage(data []byte) (*Message, error) {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, errImpossibleUnmarshalToMsg
	}

	return msg, nil
}

func MarshalToMessage(msgType string, v interface{}) (*Message, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return NewMessage(msgType, data), nil
}

func (m *Message) ToUnmarshalData(v interface{}) error {
	// Не возвращает ошибку
	data, _ := m.Data.MarshalJSON()
	if reflect.DeepEqual(data, []byte("null")) {
		return errDataIsEmpty
	}

	return json.Unmarshal(data, v)
}

func NewMessage(msgType string, data []byte) *Message {
	return &Message{
		MsgType: msgType,
		Data:    json.RawMessage(data),
	}
}
