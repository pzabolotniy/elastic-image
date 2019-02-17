package resize

import (
	"bytes"
	"github.com/nfnt/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"image"
	"image/jpeg"
	"io"
)

// Resizer is an interface that wraps Resize method
type Resizer interface {
	Logger() logging.Logger
	Resize(src io.Reader, width, height uint) (resultImage []byte, err error)
}

// Container is a container for resize parameters
type Container struct {
	logger logging.Logger
}

// NewResizer is a constructor for Resizer
func NewResizer(logger logging.Logger) Resizer {
	c := &Container{
		logger: logger,
	}
	return c
}

// Logger is a getter for Container.logger
func (c *Container) Logger() logging.Logger {
	return c.logger
}

// Resize resizes image according to width and height
func (c *Container) Resize(srcImage io.Reader, width, height uint) ([]byte, error) {
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
