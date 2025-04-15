package playing

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/jthughes/disgo/internal/repl"
)

// type Model struct {
// }

// func (m Model) Init() tea.Cmd {
// 	return nil
// }

// func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
// 	return m, nil
// }

// func (m Model) View() string {
// 	return "model"
// }

// var docStyle = lipgloss.NewStyle().Margin(1, 2)

// type Item struct {
// 	title, desc string
// }

// func NewItem(title, desc string) Item {
// 	return Item{title: title, desc: desc}
// }

// func (i Item) Title() string       { return i.title }
// func (i Item) Description() string { return i.desc }
// func (i Item) FilterValue() string { return i.title }

type Model struct {
	config *repl.Config
	// progress progress.Model
}

func New(cfg *repl.Config) Model {
	return Model{
		config: cfg,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			m.config.Player.Next()
		case "p":
			m.config.Player.Previous()
		case "space":
			if m.config.Player.Controller.Paused {
				m.config.Player.Resume()
			} else {
				m.config.Player.Pause()
			}
		case "r":
			m.config.Player.Repeat = !m.config.Player.Repeat
		}

		// case tea.WindowSizeMsg:
		// 	h, v := docStyle.GetFrameSize()
		// 	m.List.SetSize(msg.Width-h, msg.Height-v)
	}
	// var cmd tea.Cmd
	// m.List, cmd = m.List.Update(msg)
	return m, nil
}

func (m Model) View() string {
	return m.config.Player.Playlist[m.config.Player.PlaylistPosition].String()
}
