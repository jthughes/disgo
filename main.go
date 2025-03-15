package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/jthughes/gogroove/internal/database"
	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "./library.db"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Printf("failed to open db: %v", err)
		return
	}

	library := Library{
		dbq: database.New(db),
	}
	err = library.dbq.CreateTable(context.TODO())
	if err != nil {
		fmt.Printf("Failed to initialize db for library: %v\n", err)
		return
	}

	tokenOptions := policy.TokenRequestOptions{
		Scopes: []string{"User.Read", "Files.Read"},
	}

	source, err := NewOneDriveSource(tokenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	library.ImportFromSource(source, "/Music/Video Games/Darren Korb/")
}
