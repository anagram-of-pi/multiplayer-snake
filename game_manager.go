package main

import (
	"context"
	"fmt"
	"slices"
)

type PlayerID int

type Action int

const (
	JoinQueue Action = iota
	LeaveQueue
	IsReady
	ButtonPress
	ChangeName
)

// A representation of a player (?)
type NetworkingPlayer struct {
	sshID      int // ???
	PlayerName string
	PlayerID   PlayerID
}

type ClientMessage struct {
	PlayerID PlayerID
	Action   Action
}

type ServerMessage struct {
	GameState GameState
}

type Client struct {
	ID       PlayerID
	ToServer chan<- ClientMessage
	ToClient <-chan ServerMessage
}

type WaitingRoom struct {
	Players   []*NetworkingPlayer
	GameState GameState
}

type ActiveGame struct {
	Players          []*NetworkingPlayer
	PlayerDirections map[PlayerID]Direction
	GameState        GameState
}

type GameManager struct {
	NewClients <-chan Client
	ToServer   <-chan ClientMessage
}

func NewGameManager(newClients <-chan Client, toServer <-chan ClientMessage) GameManager {
	return GameManager{
		NewClients: newClients,
		ToServer:   toServer,
	}
}

func NewClient(id PlayerID, toServer chan<- ClientMessage) (Client, chan ServerMessage) {
	toClient := make(chan ServerMessage, 8)
	return Client{
		ID:       id,
		ToServer: toServer,
		ToClient: toClient,
	}, toClient
}
func (client *Client) String() string {
	return fmt.Sprintf("[Client ID #%v]", client.ID)
}

func (room *WaitingRoom) AddPlayer(player *NetworkingPlayer) {
	room.Players = append(room.Players, player)
}
func (room *WaitingRoom) RemovePlayer(player *NetworkingPlayer) {
	i := slices.Index(room.Players, player)
	room.Players = slices.Delete(room.Players, i, i+1)
}

func NewGame(directions map[PlayerID]Direction) ActiveGame {
	GAME_SIZE := 10
	return ActiveGame{
		// Players:          slices.Collect(maps.Keys(directions)),
		PlayerDirections: directions,
		GameState:        NewGameState(GAME_SIZE, GAME_SIZE),
	}
}

func (manager *GameManager) RunGameServer(ctx context.Context) {
	fmt.Println("Started running Game Server")
	for {
		select {
		case client := <-manager.NewClients:
			fmt.Printf("New Client: %s\n", client.String())
		case userEvent := <-manager.ToServer:

			fmt.Printf("New user event: %v\n", userEvent)
		case <-ctx.Done():
			fmt.Println("Stopped running Game Server")
			return
		}
	}
}
