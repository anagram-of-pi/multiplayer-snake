package main

import (
	"fmt"
	"time"
)

type GameInput struct {
	Direction Direction
}

type Point struct {
	X int
	Y int
}

// func NewPoint(x, y int) Point {...}
func P(x, y int) Point {
	return Point{x, y}
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Snake struct {
	Head      Point
	Tail      []Direction
	Direction Direction
}

type GameState struct {
	Snake  *Snake
	Apple  Point
	Speed  int
	Width  int
	Height int
}

type SnakeGame struct {
	GameState      *GameState
	commandChannel <-chan GameInput
	ticker         *time.Ticker
}

func NewGameState(width, height int) GameState {
	return GameState{
		Snake:  &Snake{},
		Apple:  Point{},
		Speed:  1,
		Width:  width,
		Height: height,
	}
}

func (p Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}
func (direction Direction) String() string {
	return [4]string{"Up", "Right", "Down", "Left"}[direction]
}
func (snake Snake) String() string {
	return fmt.Sprintf("Snake @ %s facing %s", snake.Head.String(), snake.Direction.String())
}

func (snake Snake) Length() int {
	return len(snake.Tail)
}
func (snake Snake) SetDirection(dir Direction) {
	snake.Direction = dir
}
func (snake Snake) AppendToTail(dir Direction) {
	snake.Tail = append([]Direction{dir}, snake.Tail...)
}
func (snake Snake) RemoveEndOfTail() {
	if (len(snake.Tail)) > 0 {
		snake.Tail = snake.Tail[:len(snake.Tail)-1]
	}
}

func (game *SnakeGame) RunGame() {
	ticker := game.ticker
	commandChannel := game.commandChannel
	moveDelay := time.Duration(1/game.GameState.Speed) * time.Second
	ticker.Reset(moveDelay)

	lastMove := GameInput{}
	for {
		select {
		case lastMove = <-commandChannel:
			fmt.Println(lastMove)

		case <-ticker.C:
			snake := game.GameState.Snake
			snake.AppendToTail(lastMove.Direction)
			snake.RemoveEndOfTail()
		}
	}
}

func (game *SnakeGame) ChangeSpeed(newSpeed int) {
	game.GameState.Speed = newSpeed

	moveDelay := time.Duration(1/game.GameState.Speed) * time.Second
	game.ticker.Reset(moveDelay)
}
