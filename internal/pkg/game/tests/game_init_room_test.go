package game_test

import (
	"testing"
	"time"

	g "github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/game"
)

type TestCaseInitRoom struct {
	params        g.RoomParameters
	isMustCreated bool
}

func TestGameInitRoomCreatedSameParamsTwo(t *testing.T) {
	tests := []TestCaseInitRoom{
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
	}
	game := g.NewGame()
	fID := game.InitRoom(tests[0].params)
	// Ждем пока успеет создаться 1 комната
	time.Sleep(time.Millisecond * 1)
	sID := game.InitRoom(tests[1].params)

	t.Log("First ID: ", fID)
	t.Log("Second ID: ", sID)

	if fID != sID {
		t.Fatalf("Unexpected roomid: %v. Expected: %v", sID, fID)
	}
}

func TestGameInitRoomCreatedSameParamsThree(t *testing.T) {
	tests := []TestCaseInitRoom{
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
	}
	game := g.NewGame()
	fID := game.InitRoom(tests[0].params)
	sID := game.InitRoom(tests[1].params)
	tID := game.InitRoom(tests[2].params)

	t.Log("First ID: ", fID)
	t.Log("Second ID: ", sID)
	t.Log("Third ID: ", tID)

	if fID != sID {
		t.Fatalf("Unexpected roomid: %v. Expected: %v", sID, fID)
	}

	if tID == sID {
		t.Fatalf("Unexpected result:  mustn't be equal")
	}
}

func TestGameInitRoomCreatedSameParamsFour(t *testing.T) {
	tests := []TestCaseInitRoom{
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
		TestCaseInitRoom{
			params: g.RoomParameters{time.Duration(1)},
		},
	}
	game := g.NewGame()
	fID := game.InitRoom(tests[0].params)
	// Ждем пока успеет создаться 1 комната
	time.Sleep(time.Millisecond * 1)
	sID := game.InitRoom(tests[1].params)
	tID := game.InitRoom(tests[2].params)
	// Ждем пока успеет создаться 2 комната
	time.Sleep(time.Millisecond * 1)
	frID := game.InitRoom(tests[3].params)

	t.Log("First ID: ", fID)
	t.Log("Second ID: ", sID)
	t.Log("Third ID: ", tID)
	t.Log("Fourth ID: ", frID)

	if fID != sID {
		t.Fatalf("Unexpected roomid: %v. Expected: %v", sID, fID)
	}

	if tID == sID {
		t.Fatalf("Unexpected result:  mustn't be equal")
	}

	if frID != tID {
		t.Fatalf("Unexpected roomid: %v. Expected: %v", frID, tID)
	}
}
