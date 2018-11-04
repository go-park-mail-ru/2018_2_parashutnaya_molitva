package game

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
)

const configFileName = "game.json"

type GameConfig struct {
	InitMessageDeadline int
	CloseRoomDeadline   int
	ValidGameDuration   []int
}

var (
	jsonConfigReader = config.JsonConfigReader{}
	gameConfig       = &GameConfig{}
)

func init() {
	singletoneLogger.LogMessage("INIT")
	err := jsonConfigReader.Read(configFileName, gameConfig)
	if err != nil {
		singletoneLogger.LogError(err)
	}
}
