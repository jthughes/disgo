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
	meta_has_drm,
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
    ?,
    ?
) RETURNING *;

-- name: GetTracksByAlbum :many
SELECT * FROM tracks
WHERE tracks.meta_album = ?
ORDER BY meta_track ASC;

-- name: GetAlbumsByName :many
SELECT DISTINCT meta_album FROM tracks
ORDER BY meta_album ASC;
