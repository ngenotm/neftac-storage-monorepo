package transform

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"github.com/nfnt/resize"
)

func Resize(data []byte, width uint) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil { return nil, err }
	resized := resize.Resize(width, 0, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	resize.Encode(buf, resized, &resize.ImageOptions{Format: "webp"})
	return buf.Bytes(), nil
}
