package main

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type teaModel struct {
	PlayerID        int
	PlayerName      string
	PlayerDirection Direction
	GameState       GameState
}

type PlayerMove struct {
	Direction Direction
	SenderID  int
}

type KeyMap struct {
	Up           key.Binding
	Right        key.Binding
	Down         key.Binding
	Left         key.Binding
	DirectionKey key.Binding
	Select       key.Binding
	Quit         key.Binding
}

var DefaultKeyMap = KeyMap{
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
	Quit: key.NewBinding(
		key.WithKeys("q", "escape"),
		key.WithHelp("q/esc", "quit"),
	),
}

func initialModel() teaModel {
	return teaModel{
		PlayerName: "User",
	}
}

func (model teaModel) Init() tea.Cmd {
	return nil
}

func (model teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return model, tea.Quit
		}
	}
	return model, nil
}

func (model teaModel) View() tea.View {
	return tea.NewView("")
}
