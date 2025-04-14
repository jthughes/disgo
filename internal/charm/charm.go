package charm

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/jthughes/disgo/internal/charm/albums"
	"github.com/jthughes/disgo/internal/charm/playing"
)

// https://github.com/charmbracelet/bubbletea/blob/main/examples/composable-views/main.go

type sessionState uint

const (
	defaultTime              = time.Minute
	albumsView  sessionState = iota
	playingView
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	state   sessionState
	albums  albums.Model
	playing playing.Model
}

func NewModel() model {
	items := []list.Item{
		albums.NewItem("Songs of Sanctuary", "Adiemus"),
		albums.NewItem("Sinnohvation", "insaneintherain"),
	}

	m := model{state: albumsView}
	m.albums = albums.Model{List: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.albums.List.Title = "Albums"
	m.playing = playing.Model{}
	return m
}

func (m model) Init() tea.Cmd {

	return tea.Batch(m.albums.Init(), m.playing.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == albumsView {
				m.state = playingView
			} else {
				m.state = albumsView
			}
		}
		switch m.state {
		// update whichever model is focused
		case albumsView:
			m.albums, cmd = m.albums.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.playing, cmd = m.playing.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		{
			m.albums, cmd = m.albums.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	// model := m.currentFocusedModel()
	if m.state == albumsView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", m.albums.View())), modelStyle.Render(m.playing.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.albums.View())), focusedModelStyle.Render(m.playing.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • q: exit\n"))
	return s
}

func (m model) currentFocusedModel() string {
	if m.state == albumsView {
		return "albums"
	}
	return "playing"
}
