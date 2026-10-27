package storage

import "context"

type IStorageClient interface {
	UploadFile(ctx context.Context, filename string, data []byte) (string, error)
}
