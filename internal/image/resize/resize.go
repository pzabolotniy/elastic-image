package resize

import (
	"bytes"
	"github.com/nfnt/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"image"
	"image/jpeg"
	"io"
)

type Resizer interface {
	Logger() logging.Logger
	Resize (src io.Reader, width, height uint) (resultImage []byte, err error)
}

type ResizeContainer struct {
	logger logging.Logger
}

func NewResizer( logger logging.Logger ) Resizer {
	c := &ResizeContainer{
		logger:logger,
	}
	return c
}

func (c *ResizeContainer) Logger() logging.Logger {
	return c.logger
}

func (c *ResizeContainer) Resize( srcImage io.Reader, width, height uint ) ( []byte, error ) {
	logger := c.Logger()
	buffer := new(bytes.Buffer)
	var err error

	image, _, err := image.Decode(srcImage)
	if err != nil {
		logger.Errorf("decode image failed: '%s'", err)
		return buffer.Bytes(), err
	}

	resizedImage := resize.Resize(width, height, image, resize.Lanczos3)
	err = jpeg.Encode(buffer, resizedImage, nil)
	if err != nil {
		logger.Errorf("get write for image failed: '%s'", err)
		return buffer.Bytes(), err
	}

	return buffer.Bytes(), err
}
