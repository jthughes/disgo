package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/jthughes/disgo/internal/database"
	"github.com/jthughes/disgo/internal/player"
	"github.com/jthughes/disgo/internal/player/sources"
	"github.com/jthughes/disgo/internal/repl"
	"github.com/jthughes/disgo/internal/util"
	_ "modernc.org/sqlite"
)

//go:embed sql/schema/001_tracks.sql
var schema string

func main() {
	// Set up user config folder
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

	// Setup logging
	logFile, err := os.OpenFile(fmt.Sprintf("%s/log.txt", configPath), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	dbPath := fmt.Sprintf("%s/library.db", configPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		out := fmt.Sprintf("failed to open db: %v", err)
		fmt.Println(out)
		log.Println(out)
		return
	}

	_, err = db.Exec(schema)
	if err != nil {
		out := fmt.Sprintf("Failed to initialize db for library: %v", err)
		fmt.Println(out)
		log.Println(out)
		return
	}

	library := player.InitLibrary(database.New(db), make(map[string]player.Source))

	player := player.Player{
		Playlist:         nil,
		Repeat:           false,
		PlaylistPosition: -1,
		Controller: &beep.Ctrl{
			Streamer: &util.Queue{},
		},
	}

	config := repl.InitConfig(configPath, &library, &player)

	source, err := sources.InitOneDriveSource(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	library.Sources[source.String()] = source

	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	player.Init()

	repl.Run(&config)
}
