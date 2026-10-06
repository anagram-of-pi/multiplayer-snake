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

func (game GameState) RunGame() {
	commandChannel := make(chan GameInput, 4)
	moveDelay := time.Duration(1/game.Speed) * time.Millisecond
	ticker := time.NewTicker(moveDelay)

	lastMove := GameInput{}
	for {
		select {
		case lastMove = <-commandChannel:
			fmt.Println(lastMove)

		case <-ticker.C:
			game.Snake.AppendToTail(lastMove.Direction)
			game.Snake.RemoveEndOfTail()
		}
	}
}
