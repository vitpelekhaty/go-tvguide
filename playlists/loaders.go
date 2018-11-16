package playlists

import (
	"io/ioutil"
	"net/url"
	"os"
)

// ILoader interface of playlist loaders
type ILoader interface {
	Load(path string) ([]byte, error)
}

// FileLoader - object for loading playlist data from the file
type FileLoader struct {
}

// HTTPLoader - object for downloading playlist data from the remote server
type HTTPLoader struct {
	Downloader
}

type downloadResult struct {
	data []byte
	err  error
}

// Loader return playlist loader for the specified path of the playlist
func Loader(path string) ILoader {

	if _, err := os.Stat(path); err == nil {
		return new(FileLoader)
	}

	if _, err := url.ParseRequestURI(path); err == nil {
		return new(HTTPLoader)
	}

	return nil
}

// Load returns data of the playlist with specified file path
func (loader *FileLoader) Load(path string) ([]byte, error) {

	data := make([]byte, 0)

	f, err := os.Open(path)

	if err != nil {
		return data, err
	}

	defer f.Close()

	if data, err = ioutil.ReadAll(f); err != nil {
		return data, err
	}

	return data, nil
}

// Load returns data of the playlist with specified URI
func (loader *HTTPLoader) Load(path string) ([]byte, error) {
	return loader.Run(path)
}
