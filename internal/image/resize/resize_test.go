package resize

import (
	"github.com/nfnt/resize"
	"github.com/pzabolotniy/elastic-image/internal/config"
	"github.com/pzabolotniy/elastic-image/internal/tests/mocks"
	"github.com/pzabolotniy/monkey"
	"github.com/stretchr/testify/assert"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestResizeContainer_Resize(t *testing.T) {
	testName := t.Name()
	testConfig := config.GetConfig()
	testLogger := testConfig.APILogger
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
		assert.Equalf(t, width, gotWidth, "%s - resize.Resize call, width ok", testName)
		assert.Equalf(t, heigth, gotHeight, "%s - resize.Resize call, heigth ok", testName)
		assert.Equalf(t, resize.Lanczos3, gotInterp, "%s - resize.Resize call, interpolation ok", testName)
		return mockedDstImage
	})

	monkey.Patch(jpeg.Encode, func(w io.Writer, m image.Image, o *jpeg.Options) error {
		w.Write([]byte(testImage))
		return nilError
	})

	resizer := NewResizer(testLogger)
	newImage, err := resizer.Resize(imageReader, width, heigth)

	imageBytes, err := ioutil.ReadAll(imageReader)
	if err != nil {
		panic(err)
	}
	assert.NoErrorf(t, err, "%s - no errors, ok", testName)
	assert.Equalf(t, imageBytes, newImage, "%s - result ok", testName)

	mockedSrcImage.AssertExpectations(t)
	mockedDstImage.AssertExpectations(t)
}
