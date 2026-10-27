package services

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"

	_ "golang.org/x/image/webp"
)

func TestConvertToWebP(t *testing.T) {
	// JPEG sintetis 2000x1500 — melampaui batas resize 1600px
	src := image.NewRGBA(image.Rect(0, 0, 2000, 1500))
	for y := 0; y < 1500; y++ {
		for x := 0; x < 2000; x++ {
			src.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 128, 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, src, nil); err != nil {
		t.Fatalf("encode jpeg: %v", err)
	}

	got, err := convertToWebP(buf.Bytes())
	if err != nil {
		t.Fatalf("convertToWebP: %v", err)
	}
	if len(got) < 12 || !bytes.HasPrefix(got, []byte("RIFF")) || string(got[8:12]) != "WEBP" {
		t.Fatalf("output bukan webp, header=%q", got[:min(len(got), 12)])
	}

	img, format, err := image.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatalf("decode ulang webp: %v", err)
	}
	if format != "webp" {
		t.Fatalf("format = %q, want webp", format)
	}
	b := img.Bounds()
	if longest := max(b.Dx(), b.Dy()); longest > 1600 {
		t.Fatalf("sisi terpanjang = %d, want <= 1600", longest)
	}
}
