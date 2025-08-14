package imageprocessor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/Artorison/Image-resizer/internal/models"
	"github.com/disintegration/imaging"
)

type Processor struct {
	Quality int
}

func New(quality int) *Processor {
	if quality <= 0 {
		quality = 90
	}
	if quality > 100 {
		quality = 100
	}
	fmt.Println(quality)
	return &Processor{
		Quality: quality,
	}
}

func (p *Processor) Resize(src io.Reader, width, height int) (data []byte, contentType string, err error) {
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, "", models.Wrap("decode failed", err)
	}
	processedImage := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)

	buf := new(bytes.Buffer)

	switch format {
	case "jpeg", "jpg":
		if err := jpeg.Encode(buf, processedImage, &jpeg.Options{Quality: p.Quality}); err != nil {
			return nil, "", models.Wrap("jpeg encode failed", err)
		}
		contentType = "image/jpeg"

	case "png":
		if err := png.Encode(buf, processedImage); err != nil {
			return nil, "", models.Wrap("png encode failed", err)
		}
		contentType = "image/png"

	default:
		if err := jpeg.Encode(buf, processedImage, &jpeg.Options{Quality: p.Quality}); err != nil {
			return nil, "", models.Wrap("jpeg encode failed", err)
		}
		contentType = "image/jpeg"
	}

	return buf.Bytes(), contentType, nil
}
