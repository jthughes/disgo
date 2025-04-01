package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/jthughes/disgo/internal/database"
	_ "modernc.org/sqlite"
)

//go:embed sql/schema/001_tracks.sql
var schema string

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to access user home directory for config: %v", err)
		return
	}
	configPath := fmt.Sprintf("%s/.disgo", homeDir)
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configPath, 0700)
		if err != nil {
			log.Println(err)
		}
	}

	dbPath := fmt.Sprintf("%s/library.db", configPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Printf("failed to open db: %v", err)
		return
	}

	_, err = db.Exec(schema)
	if err != nil {
		fmt.Printf("Failed to initialize db for library: %v\n", err)
		return
	}

	library := Library{
		dbq:     database.New(db),
		sources: make(map[string]Source),
	}

	player := Player{
		Playlist:         nil,
		Repeat:           false,
		PlaylistPosition: -1,
		Controller: &beep.Ctrl{
			Streamer: &Queue{},
		},
	}

	config := Config{
		configPath: configPath,
		library:    &library,
		player:     &player,
	}

	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}

	source, err := config.NewOneDriveSource(tokenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	library.sources[source.String()] = source

	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	player.Init()

	repl(&config)
}
