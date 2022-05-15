package httpclient

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/admpub/godownloader/iotools"
)

type DefaultDownloader struct {
	dp     DownloadProgress
	client http.Client
	url    string
	file   *iotools.SafeFile
}

func CreateDefaultDownloader(url string, file *iotools.SafeFile) *DefaultDownloader {
	var pd DefaultDownloader
	pd.file = file
	pd.url = url
	pd.dp.From = 0
	pd.dp.To = 1
	pd.dp.Pos = 0
	pd.dp.IsPartial = false
	return &pd
}

func (pd DefaultDownloader) GetProgress() interface{} {
	return pd.dp
}

func (pd *DefaultDownloader) BeforeRun() error {
	return nil
}

func (pd *DefaultDownloader) AfterStop() error {
	return nil
}

func (pd *DefaultDownloader) DoWork(ctx context.Context) (bool, error) {
	start := time.Now()
	//create new req
	r, err := http.NewRequestWithContext(ctx, "GET", pd.url, nil)
	if err != nil {
		return false, err
	}

	resp, err := pd.client.Do(r)
	if err != nil {
		log.Println("error: error download file", err)
		return false, err
	}
	defer resp.Body.Close()
	var written int64
	written, err = io.Copy(pd.file, resp.Body)
	if err != nil {
		return false, err
	}
	duration := time.Since(start)
	seconds := int64(duration.Seconds())
	if seconds > 0 {
		pd.dp.BytesInSecond = int64(written / seconds)
	}
	pd.dp.From = 1
	return true, nil
}

func (pd *DefaultDownloader) IsPartialDownload() bool {
	return pd.dp.IsPartial
}

func (pd *DefaultDownloader) ResetProgress() {
	pd.dp.Pos = pd.dp.From
	pd.dp.Speed = 0
}
