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

	fmt.Println("Starting up...")
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

	fmt.Println("Done!")
	fmt.Println("Authenticating with OneDrive...")
	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}

	source, err := NewOneDriveSource(tokenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Done!")

	library.sources[source.String()] = source

	repl(&library)
}
