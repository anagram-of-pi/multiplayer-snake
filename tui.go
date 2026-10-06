package main

type GameModel struct {
	PlayerID        int
	PlayerName      string
	PlayerDirection Direction
	GameState       GameState
}

type PlayerMove struct {
	Direction Direction
	SenderID  int
}
