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
