package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dbuuid "github.com/google/uuid"
	"github.com/jthughes/gogroove/internal/database"
)

type Library struct {
	dbq *database.Queries
}

func (l Library) ImportTrack(track Track) error {
	_, err := l.dbq.CreateTrack(context.TODO(), database.CreateTrackParams{
		ID:              dbuuid.New().String(),
		CreatedAt:       time.Now().String(),
		UpdatedAt:       time.Now().String(),
		FileName:        track.FileName,
		MimeType:        track.MimeType,
		MetaAlbum:       sql.NullString{String: track.Metadata.Album, Valid: true},
		MetaAlbumArtist: sql.NullString{String: track.Metadata.AlbumArtist, Valid: true},
		MetaArtist:      sql.NullString{String: track.Metadata.Artist, Valid: true},
		MetaBitrate:     sql.NullInt64{Int64: int64(track.Metadata.Bitrate), Valid: true},
		MetaDuration:    sql.NullInt64{Int64: int64(track.Metadata.Duration), Valid: true},
		MetaGenre:       sql.NullString{String: track.Metadata.Genre, Valid: true},
		MetaIsVariableBitrate: sql.NullInt64{Int64: func() int64 {
			if track.Metadata.IsVariableBitrate {
				return 1
			} else {
				return 0
			}
		}(), Valid: true},
		MetaTitle:    sql.NullString{String: track.Metadata.Title, Valid: true},
		MetaTrack:    sql.NullInt64{Int64: int64(track.Metadata.Track), Valid: true},
		MetaYear:     sql.NullInt64{Int64: int64(track.Metadata.Year), Valid: true},
		FileLocation: track.Data.GetId(),
		FileSource:   track.Data.String(),
	})
	return err
}

func (l Library) ImportFromSource(source Source, target string) {
	tracks, err := source.ScanFolder(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Found %d tracks.\nImporting...\n", len(tracks))
	tracks_imported := 0
	tracks_already_present := 0
	tracks_import_error := 0
	for _, track := range tracks {
		err = l.ImportTrack(track)
		if err != nil {
			if err.Error() == "constraint failed: UNIQUE constraint failed: tracks.file_name (2067)" {
				tracks_already_present += 1
			} else {
				tracks_import_error += 1
				fmt.Println(err)
			}
		} else {
			tracks_imported += 1
		}
	}
	fmt.Printf("Done!\n%d tracks imported, %d already present, %d failed\n", tracks_imported, tracks_already_present, tracks_import_error)

}
