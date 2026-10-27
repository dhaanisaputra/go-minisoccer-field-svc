package services

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/HugoSmits86/nativewebp"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// batas sisi terpanjang; di atas ini gambar di-resize proporsional
const maxImageSide = 1600

// convertToWebP decode gambar (jpeg/png/webp), resize kalau melebihi
// maxImageSide, lalu encode ulang sebagai webp lossless.
func convertToWebP(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if b := img.Bounds(); b.Dx() > maxImageSide || b.Dy() > maxImageSide {
		w, h := b.Dx(), b.Dy()
		nw, nh := w, h
		if w >= h {
			nw = maxImageSide
			nh = h * maxImageSide / w
			if nh < 1 {
				nh = 1
			}
		} else {
			nh = maxImageSide
			nw = w * maxImageSide / h
			if nw < 1 {
				nw = 1
			}
		}
		dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, b, draw.Src, nil)
		img = dst
	}

	var out bytes.Buffer
	if err := nativewebp.Encode(&out, img, nil); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
