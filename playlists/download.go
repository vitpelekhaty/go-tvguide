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

package playlists

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// DownloadStartEvent - an event that fires before the start of the download
type DownloadStartEvent func()

// DownloadDoneEvent - an event that fires after the download finishes
type DownloadDoneEvent func()

// DownloadProgressEvent - an event that fires on each iteration of the downloading
type DownloadProgressEvent func(complete uint64)

// Downloader that downloads the file. Notifies through events about the change of download status
type Downloader struct {
	complete uint64

	OnStart    DownloadStartEvent
	OnDone     DownloadDoneEvent
	OnProgress DownloadProgressEvent
}

func (d *Downloader) Write(data []byte) (int, error) {

	count := len(data)
	d.complete += uint64(count)

	d.progress(d.complete)

	return count, nil
}

func (d *Downloader) start() {

	if d.OnStart != nil {
		d.OnStart()
	}
}

func (d *Downloader) done() {

	if d.OnDone != nil {
		d.OnDone()
	}
}

func (d *Downloader) progress(complete uint64) {

	if d.OnProgress != nil {
		d.OnProgress(complete)
	}
}

// Run starts the download process
func (d *Downloader) Run(url string) ([]byte, error) {

	emptyData := make([]byte, 0)
	var httpClient = &http.Client{Timeout: time.Second * 10}

	resp, err := httpClient.Get(url)

	if err != nil {
		return emptyData, err
	}

	defer resp.Body.Close()

	d.start()

	data, err := ioutil.ReadAll(io.TeeReader(resp.Body, d))

	defer func() {
		d.done()
	}()

	if err != nil {
		return emptyData, err
	}

	return data, nil
}
