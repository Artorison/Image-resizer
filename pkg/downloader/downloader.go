package downloader

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Artorison/Image-resizer/pkg/logger"
)

type Downloader struct {
	Client *http.Client
	Logg   logger.Logger
}

func New(client *http.Client, logg logger.Logger) *Downloader {
	return &Downloader{
		Client: client,
		Logg:   logg,
	}
}

func (d *Downloader) Get(ctx context.Context, url string, headers http.Header) (io.ReadCloser, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header = headers.Clone()
	start := time.Now()
	resp, err := d.Client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return nil, 0, err
	}
	d.Logg.Info("response", slog.Any("duration", duration), slog.String("status", resp.Status))
	return resp.Body, resp.StatusCode, nil
}
