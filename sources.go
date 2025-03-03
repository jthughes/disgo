package main

import (
	"io"
)

type FileSource interface {
	Get() (io.ReadCloser, error)
}

type OneDriveFileSource struct {
	id     string
	source *OneDriveSource
}

func (file OneDriveFileSource) Get() (io.ReadCloser, error) {
	return file.source.DownloadFile(file.id)
}

// type DropboxFileSource struct {
// 	id     string
// 	source *DropboxSource
// }

// func (file DropboxFileSource) Get() (io.ReadCloser, error) {
// 	return file.source.DownloadFile(file.id)
// }
