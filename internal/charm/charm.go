package charm

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/jthughes/disgo/internal/charm/albums"
	"github.com/jthughes/disgo/internal/charm/playing"
	"github.com/jthughes/disgo/internal/charm/playlist"
	"github.com/jthughes/disgo/internal/player"
	"github.com/jthughes/disgo/internal/repl"
)

// https://github.com/charmbracelet/bubbletea/blob/main/examples/composable-views/main.go

type sessionState uint

const (
	albumsView sessionState = iota
	playingView
	playlistView
)

type model struct {
	cfg      *repl.Config
	state    sessionState
	albums   albums.Model
	playing  playing.Model
	playlist playlist.Model
	width    int
	height   int
}

func NewModel(cfg *repl.Config) model {
	allAlbums, _ := cfg.Library.GetAlbums()

	items := make([]list.Item, len(allAlbums))
	for i, album := range allAlbums {
		items[i] = albums.NewItem(
			album.Title,
			fmt.Sprintf("%s (%s) - %d tracks", album.Artist, album.Year, len(album.Tracks)))
	}

	albumsModel := albums.Model{
		List: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	albumsModel.List.Title = "Albums"
	albumsModel.List.DisableQuitKeybindings()
	return model{
		cfg:    cfg,
		state:  albumsView,
		albums: albumsModel,
	}
}

func (m model) Init() tea.Cmd {
	return nil
	// return tea.Batch(m.albums.Init(), m.playing.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			{
				switch m.state {
				case albumsView:
					{
						// Get selection
						selectedAlbumIndex := m.albums.List.GlobalIndex()

						// Get album from selection
						albums, _ := m.cfg.Library.GetAlbums()
						album := albums[selectedAlbumIndex]
						// Add album to playlist
						if len(m.cfg.Player.Playlist) > 0 {
							m.cfg.Player.Stop()
						}
						m.cfg.Player.AddAlbumToPlaylist(album)

						// Create playlist model
						m.playlist = newPlaylistModel(m.cfg.Player.Playlist)
						m.playlist.List.DisableQuitKeybindings()

						msg := tea.WindowSizeMsg{
							Width:  m.width,
							Height: m.height,
						}
						m.playlist, _ = m.playlist.Update(msg)

						// Create playing model
						m.playing = playing.New(m.cfg)

						m.cfg.Player.Play()

						// Change to Now Playing
						m.state = playingView
					}
				case playlistView:
					{
						selectedAlbumIndex := m.playlist.List.GlobalIndex()
						offset := selectedAlbumIndex - m.cfg.Player.PlaylistPosition
						m.cfg.Player.JumpTo(offset)
						m.state = playingView
					}
				}
			}
		case "esc":
			if m.state != albumsView {
				m.state = albumsView
			}
		case "tab":
			{
				switch m.state {
				case playingView:
					{
						m.state = playlistView
					}
				case playlistView:
					{
						m.state = playingView
					}
				}
			}
		case "s":
			m.cfg.Library.ImportFromSource(m.cfg.Library.Sources["onedrive"], "/Music/")
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	switch m.state {
	case albumsView:
		m.albums, cmd = m.albums.Update(msg)
		cmds = append(cmds, cmd)
	case playingView:
		m.playing, cmd = m.playing.Update(msg)
		cmds = append(cmds, cmd)
	case playlistView:
		m.playlist, cmd = m.playlist.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
	if m.state == albumsView {
		s += m.albums.View()
	} else if m.state == playingView {
		s += m.playing.View()
	} else if m.state == playlistView {
		s += m.playlist.View()
	}
	return s
}

func newPlaylistModel(tracks []player.Track) playlist.Model {
	items := make([]list.Item, len(tracks))
	for i, track := range tracks {
		trackDuration := time.Duration(track.Metadata.Duration * int(time.Millisecond)).Truncate(time.Second).String()
		title := fmt.Sprintf("%d. %s", track.Metadata.Track, track.Metadata.Title)
		log.Println(title)
		items[i] = playlist.NewItem(
			// fmt.Sprintf("%d. %s", track.Metadata.Track, track.Metadata.Title),
			title,
			fmt.Sprintf("%s (%s)", track.Metadata.Artist, trackDuration))
	}
	model := playlist.Model{
		List: list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	model.List.Title = "Playlist"
	return model
}
