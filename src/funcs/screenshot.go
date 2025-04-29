package funcs

import (
	"bytes"
	"image/jpeg"

	"github.com/vova616/screenshot"
)

func Screenshot() (*bytes.Buffer, int) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return nil, -2
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, -1
	}

	return &buf, 0
}
