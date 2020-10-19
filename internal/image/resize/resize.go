package resize

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"

	"github.com/pzabolotniy/elastic-image/internal/logging"

	"github.com/nfnt/resize"
)

// Resize resizes image according to width and height
func Resize(ctx context.Context, srcImage io.Reader, width, height uint) ([]byte, error) {
	logger := logging.FromContext(ctx)
	buffer := new(bytes.Buffer)

	decodedImage, _, err := image.Decode(srcImage)
	if err != nil {
		logger.WithError(err).Error("decode image failed")
		return nil, err
	}

	resizedImage := resize.Resize(width, height, decodedImage, resize.Lanczos3)
	err = jpeg.Encode(buffer, resizedImage, nil)
	if err != nil {
		logger.WithError(err).Error("jpeg Encode image failed")
		return nil, err
	}

	return buffer.Bytes(), err
}
