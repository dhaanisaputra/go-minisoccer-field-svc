package services

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"testing"

	_ "golang.org/x/image/webp"
)

func TestConvertToWebPRejectsHugeDimensions(t *testing.T) {
	// header JPEG minimal yang mendeklarasikan 10000x10000 (tanpa data scan)
	// header JPEG minimal: SOI + SOF0 (deklarasi dimensi) + SOS + EOI —
	// cukup buat DecodeConfig, tanpa data scan
	data := []byte{0xFF, 0xD8, 0xFF, 0xC0, 0x00, 0x11, 0x08,
		0x27, 0x11, // height 10001 (lewat batas 10000)
		0x27, 0x11, // width 10001
		0x03,
		0x01, 0x11, 0x00,
		0x02, 0x11, 0x00,
		0x03, 0x11, 0x00,
		0xFF, 0xDA, 0x00, 0x0C, 0x03, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x00, 0x3F, 0x00,
		0xFF, 0xD9}

	_, err := convertToWebP(data)
	if err == nil {
		t.Fatal("expected error untuk gambar 10001x10001")
	}
	if !errors.Is(err, ErrImageTooLarge) {
		t.Fatalf("err = %v, want ErrImageTooLarge", err)
	}
}

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
