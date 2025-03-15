-- name: CreateTable :exec
CREATE TABLE IF NOT EXISTS tracks (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	file_name TEXT UNIQUE NOT NULL,
    mime_type TEXT NOT NULL,
    meta_album TEXT,
	meta_album_artist TEXT,
	meta_artist TEXT,
	meta_bitrate INTEGER,
	meta_duration INTEGER,
	meta_genre TEXT,
	meta_is_variable_bitrate INTEGER,
	meta_title TEXT,
	meta_track INTEGER,
	meta_year INTEGER,
	file_location TEXT NOT NULL,
	file_source TEXT NOT NULL
);

-- name: CreateTrack :one
INSERT INTO tracks (
    id,
	created_at,
    updated_at,
	file_name,
    mime_type,
    meta_album,
	meta_album_artist,
	meta_artist,
	meta_bitrate,
	meta_duration,
	meta_genre,
	meta_is_variable_bitrate,
	meta_title,
	meta_track,
	meta_year,
	file_location,
	file_source)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
) RETURNING *;

-- name: GetTracksByAlbum :many
SELECT * FROM tracks
WHERE tracks.meta_album = ?
ORDER BY meta_track ASC;

-- name: GetAlbumsByName :many
SELECT DISTINCT meta_album FROM tracks
ORDER BY meta_album DESC;
