package main

import "fmt"

type GameInput struct {
	Direction Direction
}

type Point struct {
	X int
	Y int
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

func SnakeGame() {
	commandChannel := make(chan GameInput, 4)
	for {
		move := <-commandChannel
		fmt.Println(move)
	}
}
