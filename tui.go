package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type teaModel struct {
	keys            keyMap
	help            help.Model
	termHeight      int
	PlayerID        PlayerID
	PlayerName      string
	PlayerDirection Direction
	Client          Client
	GameState       GameState
}

type PlayerMove struct {
	Direction Direction
	SenderID  int
}

type keyMap struct {
	Up           key.Binding
	Right        key.Binding
	Down         key.Binding
	Left         key.Binding
	DirectionKey key.Binding
	Select       key.Binding
	Help         key.Binding
	Quit         key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("w", "up"),   // actual keybindings
		key.WithHelp("↑/w", "up"), // corresponding help text
	),
	Right: key.NewBinding(
		key.WithKeys("d", "right"),
		key.WithHelp("→/d", "right"),
	),
	Down: key.NewBinding(
		key.WithKeys("s", "down"),
		key.WithHelp("↓/s", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("a", "left"),
		key.WithHelp("←/a", "left"),
	),
	DirectionKey: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "use button"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.DirectionKey, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.DirectionKey}, // first column
		{k.Select, k.Help, k.Quit},                      // second column
	}
}

func initialModel(client Client) teaModel {
	return teaModel{
		PlayerName: "User",
		Client:     client,
		keys:       keys,
		help:       help.New(),
	}
}

func (model teaModel) Init() tea.Cmd {
	return nil
}

func (model teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model.GameState.Width = 10
	model.GameState.Height = 10
	model.GameState.Snake = &Snake{
		Tail: []Direction{Up, Up, Left},
	}

	select {
	case msg := <-model.Client.ToClient:
		if (msg.GameState != GameState{}) {
			model.GameState = msg.GameState
		}

	default:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			model.help.SetWidth(msg.Width)
			model.termHeight = msg.Height
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, keys.Quit):
				toServer := model.Client.ToServer
				if toServer != nil {
					toServer <- ClientMessage{
						model.PlayerID,
						LeaveQueue,
					}
				}
				return model, tea.Quit
			}
		}
	}

	return model, nil
}

func (model teaModel) View() tea.View {
	board := renderBoard(model.GameState)

	helpView := model.help.View(model.keys)
	height := (model.termHeight - 2) - strings.Count(board, "\n") - strings.Count(helpView, "\n")
	height = max(height, 0)

	view := tea.NewView(board + strings.Repeat("\n", height) + helpView)
	view.AltScreen = true
	return view
}

func renderBoard(gameState GameState) string {
	if gameState.Width == 0 || gameState.Height == 0 || gameState.Snake == nil {
		return "Error: Invalid board."
	}

	board := make([][]string, gameState.Height)
	for i := range gameState.Height {
		row := make([]string, gameState.Width)
		for i := range row {
			row[i] = " "
		}
		board[i] = row
	}

	apple := gameState.Apple
	board[apple.Y][apple.X] = "@"
	currentSnakePos := gameState.Snake.Head
	for i := range gameState.Snake.Tail {
		board[currentSnakePos.Y][currentSnakePos.X] = "#"
		currentSnakePos = gameState.Snake.Tail[i].AddTo(currentSnakePos)
	}
	headPos := gameState.Snake.Head
	board[headPos.Y][headPos.X] = "O"

	var sb strings.Builder
	sb.Grow(gameState.Width * gameState.Height * 4)

	for _, row := range board {
		for _, cell := range row {
			fmt.Fprintf(&sb, "[%s]", cell)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func tuiMain(client Client) {
	p := tea.NewProgram(initialModel(client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
