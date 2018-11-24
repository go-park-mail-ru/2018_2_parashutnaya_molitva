package game

import (
	"encoding/json"

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

//easyjson:json
type ResultMessage struct {
	Result string `json:"result"`
	Score  int    `json:"score"`
}

//easyjson:json
type TurnMessage struct {
	Turn string `json:"turn"`
}

//easyjson:json
type StartMessage struct {
	Color string `json:"color"`
	GUID  string `json:"guid"`
}

//easyjson:json
type InitMessage struct {
	RoomId string `json:"roomid"`
}

//easyjson:json
type ErrorMessage struct {
	Error string `json:"error"`
}

//easyjson:json
type InfoMessage struct {
	Info string `json:"info"`
}

//easyjson:json
type Message struct {
	MsgType string          `json:"MsgType"`
	Data    json.RawMessage `json:"Data"`
}

func UnmarshalToMessage(data []byte) (*Message, error) {
	msg := &Message{}
	err := msg.UnmarshalJSON(data)
	if err != nil {
		return nil, errImpossibleUnmarshalToMsg
	}

	return msg, nil
}

func MarshalToMessage(msgType string, v json.Marshaler) (*Message, error) {

	data, err := v.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return NewMessage(msgType, data), nil
}

func (m *Message) UnmarshalData(v json.Unmarshaler) error {
	if m.Data == nil {
		return errDataIsEmpty
	}
	data := []byte(m.Data)
	return v.UnmarshalJSON(data)
}

func NewMessage(msgType string, data []byte) *Message {
	return &Message{
		MsgType: msgType,
		Data:    json.RawMessage(data),
	}
}
