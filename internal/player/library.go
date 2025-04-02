package player

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dbuuid "github.com/google/uuid"
	"github.com/jthughes/disgo/internal/database"
)

type Library struct {
	dbq     *database.Queries
	Sources map[string]Source
}

func InitLibrary(dbq *database.Queries, src map[string]Source) Library {
	return Library{
		dbq:     dbq,
		Sources: src,
	}
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
		FileLocation: track.Data.location,
		FileSource:   track.Data.sourceName,
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

func (l Library) GetAlbums() ([]Album, error) {
	albums := make([]Album, 0)
	dbAlbums, err := l.dbq.GetAlbumsByName(context.TODO())
	if err != nil {
		return nil, err
	}
	for _, album := range dbAlbums {
		if !album.Valid {
			continue
		}
		dbTracks, err := l.dbq.GetTracksByAlbum(context.TODO(), album)
		if err != nil {
			return nil, err
		}
		if len(dbTracks) == 0 {
			continue
		}

		tracks := make([]Track, 0)
		for _, track := range dbTracks {

			tracks = append(tracks, Track{
				Data: File{
					location:   track.FileLocation,
					sourceName: track.FileSource,
					source:     l.Sources[track.FileSource],
				},
				FileName: track.FileName,
				Metadata: AudioMetadata{
					Album:             fromDBString(track.MetaAlbum),
					AlbumArtist:       fromDBString(track.MetaAlbumArtist),
					Artist:            fromDBString(track.MetaArtist),
					Bitrate:           fromDBInt(track.MetaBitrate),
					Duration:          fromDBInt(track.MetaDuration),
					Genre:             fromDBString(track.MetaGenre),
					HasDrm:            fromDBBool(track.MetaHasDrm),
					IsVariableBitrate: fromDBBool(track.MetaIsVariableBitrate),
					Title:             fromDBString(track.MetaTitle),
					Track:             fromDBInt(track.MetaTrack),
					Year:              fromDBInt(track.MetaYear),
				},
				MimeType: track.MimeType,
			})

		}

		yearString := ""
		year := tracks[0].Metadata.Year
		if year > 0 {
			yearString = fmt.Sprint(year)
		}
		albums = append(albums, Album{
			Title:  album.String,
			Artist: fromDBString(dbTracks[0].MetaAlbumArtist),
			Year:   yearString,
			Genre:  fromDBString(dbTracks[0].MetaGenre),
			Tracks: tracks,
		})
	}
	return albums, nil
}

func fromDBString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
func fromDBInt(ns sql.NullInt64) int {
	if ns.Valid {
		return int(ns.Int64)
	}
	return 0
}

func fromDBBool(ns sql.NullInt64) bool {
	if ns.Valid && ns.Int64 == 1 {
		return true
	}
	return false
}
