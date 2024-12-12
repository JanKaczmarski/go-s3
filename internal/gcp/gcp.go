package gcp

import (
	"context"
	"time"

	mytypes "github.com/jankaczmarski/go-s3/pkg/types"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GcpBucket struct {
	bucketID  string
	projectID string
}
type GcpWorker struct {
	client    *storage.Client
	projectID string
}

func NewWorker(ctx context.Context, projectID string) (*GcpWorker, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GcpWorker{
		client:    client,
		projectID: projectID,
	}, nil
}

func NewGcpBucket(ctx context.Context, bucketID string, worker *GcpWorker) (*GcpBucket, error) {
	client := worker.client
	projectID := worker.projectID
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketID)

	if err := bucket.Create(ctx, projectID, nil); err != nil {
		return nil, err
	}

	return &GcpBucket{
		bucketID:  bucketID,
		projectID: worker.projectID,
	}, nil
}

func (worker *GcpWorker) Close() error {
	err := worker.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (worker *GcpWorker) CreateBucket(ctx context.Context, bucketID string) (*mytypes.Bucket, error) {
	gcpBucket, err := NewGcpBucket(ctx, bucketID, worker)
	if err != nil {
		return nil, err
	}
	return &mytypes.Bucket{
		Name: gcpBucket.bucketID,
	}, nil
}

func (worker *GcpWorker) ListBuckets(ctx context.Context) ([]string, error) {
	client := worker.client
	projectID := worker.projectID

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	buckets := make([]string, 0)
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}
	return buckets, nil
}
