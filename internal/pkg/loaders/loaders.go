// IPTV guide viewer
//
// Copyright 2018 Vitaly Pelekhaty
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the License.

package loaders

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
