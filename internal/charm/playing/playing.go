package playing

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Model struct {
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "model"
}
