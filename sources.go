package main

import "io"

type Source interface {
	ScanFolder(string) ([]Track, error)
	DownloadFile(string) (io.ReadCloser, error)
}

type File struct {
	location   string
	sourceName string
	source     Source
}

func (f File) Get() (io.ReadCloser, error) {
	return f.source.DownloadFile(f.location)
}

// type FileSource interface {
// 	GetId() string
// 	String() string
// 	Get() (io.ReadCloser, error)
// }

// type OneDriveFileSource struct {
// 	id     string
// 	source *OneDriveSource
// }

// func (file OneDriveFileSource) Get() (io.ReadCloser, error) {
// 	return file.source.DownloadFile(file.id)
// }

// func (file OneDriveFileSource) GetId() string {
// 	return file.id
// }

// func (file OneDriveFileSource) String() string {
// 	return "onedrive"
// }

// type DropboxFileSource struct {
// 	id     string
// 	source *DropboxSource
// }

// func (file DropboxFileSource) Get() (io.ReadCloser, error) {
// 	return file.source.DownloadFile(file.id)
// }
