package game

import (
	"fmt"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"github.com/gorilla/websocket"
)

type PlayerData struct {
	Name  string
	Score int
}

type Player struct {
	conn       *websocket.Conn
	playerData PlayerData
}

func NewPlayer(name string, score int, conn *websocket.Conn) *Player {
	return &Player{
		conn: conn,
		playerData: PlayerData{
			name,
			score,
		},
	}
}

// Читает сообщения из соединения и записывает в канал
func (p *Player) read(out chan<- *Message) {

}

// Читает сообщения из канала и записывает в соединение
func (p *Player) write(in <-chan *Message) {

}

func (p *Player) Start(out chan<- *Message, in <-chan *Message) {
	singletoneLogger.LogMessage(fmt.Sprintf("Start listen player: %v", p.playerData.Name))
	go p.read(out)
	go p.write(in)
}
