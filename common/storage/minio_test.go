package storage

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"field-service/config"
)

func minioReachable() bool {
	conn, err := net.DialTimeout("tcp", "localhost:9100", time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func TestUploadFile(t *testing.T) {
	if !minioReachable() {
		t.Skip("minio tidak jalan di localhost:9100")
	}

	cfg := config.MinIOConfig{
		Endpoint:   "localhost:9100",
		PublicURL:  "http://localhost:9100",
		AccessKey:  "minioadmin",
		SecretKey:  "minioadmin",
		UseSSL:     false,
		BucketName: "field-images",
	}
	client, err := NewMinIOClient(cfg)
	if err != nil {
		t.Fatalf("NewMinIOClient: %v", err)
	}

	ctx := context.Background()
	name := "images/test-" + time.Now().Format("20060102150405") + ".webp"
	want := []byte("RIFF\x00\x00\x00\x00WEBPtest-body")

	url, err := client.UploadFile(ctx, name, want)
	if err != nil {
		t.Fatalf("UploadFile: %v", err)
	}
	wantURL := "http://localhost:9100/field-images/" + name
	if url != wantURL {
		t.Fatalf("url = %q, want %q", url, wantURL)
	}

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want 200", url, resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("body = %q, want %q", got, want)
	}
}
