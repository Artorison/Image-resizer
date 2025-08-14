package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Artorison/Image-resizer/internal/models"
	"github.com/Artorison/Image-resizer/pkg/helpers"
	"github.com/Artorison/Image-resizer/pkg/logger"
	"github.com/labstack/echo/v4"
)

type ImageService interface {
	GetImage(ctx context.Context, params *models.ImageParams, headers http.Header) (*models.ImageResult, error)
	CheckInCache(params *models.ImageParams) (*models.ImageResult, error)
}
type Handlers struct {
	S    ImageService
	Logg logger.Logger
}

func New(service ImageService, logg logger.Logger) *Handlers {
	return &Handlers{
		S:    service,
		Logg: logg,
	}
}

func (h *Handlers) GetImage(c echo.Context) error {
	imageParams, err := parseParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Err(err.Error()))
	}

	res, err := h.S.CheckInCache(imageParams)
	if err == nil {
		c.Response().Header().Set("X-Cache", "HIT")
		return c.Blob(http.StatusOK, res.ContentType, res.Data)
	}
	if !errors.Is(err, models.ErrNoInCache) {
		h.Logg.Error("INTERNAL", slog.String("ERROR", err.Error()))
		return c.JSON(http.StatusInternalServerError, models.Err("internal"))
	}

	headers := c.Request().Header.Clone()
	headers.Set("Accept-Encoding", "identity")
	res, err = h.S.GetImage(c.Request().Context(), imageParams, headers)
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Err(err.Error()))
	}
	c.Response().Header().Set("X-Cache", "MISS")
	return c.Blob(http.StatusOK, res.ContentType, res.Data)
}

func parseParams(c echo.Context) (*models.ImageParams, error) {
	wStr, hStr := c.Param("width"), c.Param("height")

	w, err := strconv.Atoi(wStr)
	if err != nil || w <= 0 {
		return nil, fmt.Errorf("invalid width: %v, err: %w", wStr, err)
	}
	h, err := strconv.Atoi(hStr)
	if err != nil || h <= 0 {
		return nil, fmt.Errorf("invalid height: %v, err: %w", hStr, err)
	}

	raw := strings.TrimPrefix(c.Param("*"), "/")
	if raw == "" {
		return nil, fmt.Errorf("empty source URL")
	}

	if q := c.Request().URL.RawQuery; q != "" {
		raw += "?" + q
	}

	var full string
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		full = raw
	} else {
		full = "https://" + raw
	}

	srcURL, err := url.Parse(full)
	if err != nil {
		return nil, fmt.Errorf("invalid source URL %v: %w", full, err)
	}

	key := helpers.GenerateCacheKey(w, h, srcURL.String())

	return &models.ImageParams{
		Width:      w,
		Height:     h,
		SourceURL:  srcURL,
		OriginLink: full,
		CacheKey:   key,
	}, nil
}
