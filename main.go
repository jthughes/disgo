package main

import (
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/jthughes/gogroove/internal/database"
	_ "modernc.org/sqlite"
)

//go:embed sql/schema/001_tracks.sql
var schema string

func main() {
	dbPath := "./library.db"
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

	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}

	source, err := NewOneDriveSource(tokenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	library.sources[source.String()] = source

	library.ImportFromSource(source, "/Music/Video Games/Darren Korb/")
	albums, err := library.GetAlbums()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, album := range albums {
		fmt.Printf("[%d] %s (%s): %d tracks\n", i, album.Title, album.Year, len(album.Tracks))
	}
	album := albums[1]
	album.Tracks[2].Play()
}
