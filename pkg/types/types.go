package types

import (
	"context"
)

type Bucket struct {
	Name   string
	Region string
}

// worker for working with storage resources (s3 like)
type StorageWorker interface {
	ListBuckets(context.Context) ([]string, error)
	CreateBucket(ctx context.Context, name string, region string) (*Bucket, error)
}
