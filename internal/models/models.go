package models

import "net/url"

type ImageParams struct {
	Width      int
	Height     int
	OriginLink string
	SourceURL  *url.URL
	CacheKey   string
}

type ImageResult struct {
	Data        []byte
	ContentType string
}
