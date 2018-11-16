package playlists

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// DownloadStartEvent - an event that fires before the start of the download
type DownloadStartEvent func(total uint64)

// DownloadDoneEvent - an event that fires after the download finishes
type DownloadDoneEvent func()

// DownloadProgressEvent - an event that fires on each iteration of the downloading
type DownloadProgressEvent func(complete, total uint64)

// Downloader that downloads the file. Notifies through events about the change of download status
type Downloader struct {
	Total, Complete uint64

	OnStart    DownloadStartEvent
	OnDone     DownloadDoneEvent
	OnProgress DownloadProgressEvent
}

func (d *Downloader) Write(data []byte) (int, error) {

	count := len(data)
	d.Complete += uint64(count)

	d.progress(d.Complete, d.Total)

	return count, nil
}

func (d *Downloader) start(total uint64) {

	d.Total = total

	if d.OnStart != nil {
		d.OnStart(total)
	}
}

func (d *Downloader) done() {

	if d.OnDone != nil {
		d.OnDone()
	}
}

func (d *Downloader) progress(complete, total uint64) {

	if d.OnProgress != nil {
		d.OnProgress(complete, total)
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

	fsize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	d.start(uint64(fsize))

	data, err := ioutil.ReadAll(io.TeeReader(resp.Body, d))

	defer func() {
		d.done()
	}()

	if err != nil {
		return emptyData, err
	}

	return data, nil
}
