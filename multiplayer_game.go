package main

import (
	"maps"
	"slices"
)

type NetworkingPlayer struct {
	sshID      int // ???
	PlayerName string
	PlayerID   int
}

type WaitingRoom struct {
	Players   []*NetworkingPlayer
	GameState GameState
}

type ActiveGame struct {
	Players          []*NetworkingPlayer
	PlayerDirections map[*NetworkingPlayer]Direction
	GameState        GameState
}

func (room WaitingRoom) AddPlayer(player *NetworkingPlayer) {
	room.Players = append(room.Players, player)
}
func (room WaitingRoom) RemovePlayer(player *NetworkingPlayer) {
	i := slices.Index(room.Players, player)
	room.Players = slices.Delete(room.Players, i, i+1)
}

func NewGame(directions map[*NetworkingPlayer]Direction) ActiveGame {
	GAME_SIZE := 10
	return ActiveGame{
		Players:          slices.Collect(maps.Keys(directions)),
		PlayerDirections: directions,
		GameState:        NewGameState(GAME_SIZE, GAME_SIZE),
	}
}
