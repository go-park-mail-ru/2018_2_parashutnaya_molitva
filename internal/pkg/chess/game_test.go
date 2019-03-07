package chess

import (
	"fmt"
	"testing"
)

func TestGame_BoardString(t *testing.T) {
	game := NewGame()
	boardString := game.BoardString()
	expected := "RNBQKBNRPPPPPPPP................................pppppppprnbqkbnr"
	if boardString != expected {
		t.Error(fmt.Sprintf("expected: %s\ngot: %s\n", expected, boardString))
	}
}

func TestGame_CurrentTurn(t *testing.T) {
	game := NewGame()
	turn := game.CurrentTurn()
	expected := true
	if turn != expected {
		t.Error(fmt.Sprintf("expected: %t\ngot: %t\n", expected, turn))
	}
}

func TestGame_Move(t *testing.T) {
	game := NewGame()
	game.Move("e2e4")
	position := game.BoardString()
	expectedPosition := "RNBQKBNRPPPP.PPP....E.......P...................pppppppprnbqkbnr"
	if position != expectedPosition {
		t.Error(fmt.Sprintf("expected: %s\ngot: %s\n", position, expectedPosition))
	}
}

func TestGame_Status(t *testing.T) {
	game := NewGame()
	status := game.Status()
	expected := InProgress
	if status != expected {
		t.Error(fmt.Sprintf("expected: %#v\ngot: %#v\n", expected, status))
	}
}

func TestGame_StatusWhiteWon(t *testing.T) {
	game := NewGame()
	game.Move("e2e4")
	game.Move("f7f5")
	game.Move("e4f5")
	game.Move("g7g5")
	game.Move("d1h5")
	status := game.Status()
	expected := WhiteWon
	if status != expected {
		t.Error(fmt.Sprintf("expected: %#v\ngot: %#v\n", expected, status))
	}
}

func TestGame_StatusBlackWon(t *testing.T) {
	game := NewGame()
	game.Move("f2f3")
	game.Move("e7e5")
	game.Move("g2g4")
	game.Move("d8h4")
	status := game.Status()
	expected := BlackWon
	if status != expected {
		t.Error(fmt.Sprintf("expected: %#v\ngot: %#v\n", expected, status))
	}
}

func TestGame_IsCheckmate(t *testing.T) {
	game := NewGame()
	game.Move("f2f3")
	game.Move("e7e5")
	game.Move("g2g4")
	game.Move("d8h4")
	status := game.IsCheckmate()
	expected := true
	if status != expected {
		t.Error(fmt.Sprintf("expected: %#v\ngot: %#v\n", expected, status))
	}
}

func TestGame_IsCheckmate2(t *testing.T) {
	game := NewGame()
	status := game.IsCheckmate()
	expected := false
	if status != expected {
		t.Error(fmt.Sprintf("expected: %#v\ngot: %#v\n", expected, status))
	}
}

func TestGame_IsGameOver(t *testing.T) {
	game := NewGame()
	game.Move("f2f3")
	game.Move("e7e5")
	game.Move("g2g4")
	game.Move("d8h4")
	status := game.IsGameOver()
	expected := true
	if status != expected {
		t.Error(fmt.Sprintf("expected: %#v\ngot: %#v\n", expected, status))
	}
}

