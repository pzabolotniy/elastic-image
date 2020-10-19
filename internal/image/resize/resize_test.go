package resize

import (
	"context"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nfnt/resize"
	"github.com/pzabolotniy/elastic-image/internal/logging"
	"github.com/pzabolotniy/elastic-image/internal/tests/mocks"
	"github.com/pzabolotniy/monkey"
	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	ctx := context.Background()
	testLogger := logging.GetLogger()
	ctx = logging.WithContext(ctx, testLogger)
	var nilError error

	testImage := "i.am.image"
	imageReader := strings.NewReader(testImage)
	width := uint(1024)
	heigth := uint(800)

	mockedSrcImage := &mocks.Image{}
	mockedImageFormat := "jpeg"
	monkey.Patch(image.Decode, func(src io.Reader) (image.Image, string, error) {
		return mockedSrcImage, mockedImageFormat, nilError
	})

	mockedDstImage := &mocks.Image{}
	monkey.Patch(resize.Resize, func(gotWidth, gotHeight uint, img image.Image, gotInterp resize.InterpolationFunction) image.Image {
		assert.Equal(t, width, gotWidth, "resize.Resize call, width ok")
		assert.Equal(t, heigth, gotHeight, "resize.Resize call, heigth ok")
		assert.Equal(t, resize.Lanczos3, gotInterp, "resize.Resize call, interpolation ok")
		return mockedDstImage
	})

	monkey.Patch(jpeg.Encode, func(w io.Writer, m image.Image, o *jpeg.Options) error {
		w.Write([]byte(testImage))
		return nilError
	})

	newImage, err := Resize(ctx, imageReader, width, heigth)

	imageBytes, err := ioutil.ReadAll(imageReader)
	if err != nil {
		panic(err)
	}
	assert.NoError(t, err, "no errors, ok")
	assert.Equal(t, imageBytes, newImage, "result ok")

	mockedSrcImage.AssertExpectations(t)
	mockedDstImage.AssertExpectations(t)
}
