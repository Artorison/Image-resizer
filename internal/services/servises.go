package services

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Artorison/Image-resizer/internal/models"
	lrucache "github.com/Artorison/Image-resizer/pkg/lru_cache"
)

type Service struct {
	Cache      lrucache.Cache
	Processor  Processor
	Downloader Downloader
}

type Processor interface {
	Resize(src io.Reader, width, height int) (data []byte, contentType string, err error)
}

type Downloader interface {
	Get(ctx context.Context, url string, headers http.Header) (io.ReadCloser, int, error)
}

func New(cache lrucache.Cache, downloader Downloader, processor Processor) *Service {
	return &Service{
		Cache:      cache,
		Downloader: downloader,
		Processor:  processor,
	}
}

func (s *Service) GetImage(
	ctx context.Context,
	params *models.ImageParams,
	headers http.Header) (
	*models.ImageResult, error,
) {
	body, status, err := s.Downloader.Get(ctx, params.SourceURL.String(), headers)
	if err != nil {
		return nil, models.Wrap("download failed", err)
	}
	defer body.Close()

	if status >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("downloader returned status %d", status)
	}

	resized, contentType, err := s.Processor.Resize(body, params.Width, params.Height)
	if err != nil {
		return nil, models.Wrap("resize failed", err)
	}
	imageResult := &models.ImageResult{Data: resized, ContentType: contentType}
	s.Cache.Set(params.CacheKey, imageResult)

	return imageResult, nil
}

func (s *Service) CheckInCache(
	params *models.ImageParams) (
	*models.ImageResult, error,
) {
	if data, ok := s.Cache.Get(params.CacheKey); ok {
		bData, ok := data.(*models.ImageResult)
		if !ok {
			return nil, fmt.Errorf("data is broken")
		}
		return bData, nil
	}
	return nil, models.ErrNoInCache
}
