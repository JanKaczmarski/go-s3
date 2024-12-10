package gcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

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

func (worker *GcpWorker) Close() error {
	err := worker.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (worker *GcpWorker) Run() {
	// context for working with API's, sender sends context and receiver receieves it,
	// imo it's a context of certain api calls, more info here: https://pkg.go.dev/context
	ctx := context.Background()

	buckets, err := worker.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("Failed to listBuckets: %v", err)
	}

	for i, tmpBucketName := range buckets {
		fmt.Printf("ID-%d: BucketName: %s\n", i, tmpBucketName)
	}
}

func (worker *GcpWorker) CreateBucket(bucketID string, ctx context.Context) error {
	client := worker.client
	projectID := worker.projectID
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketID)

	if err := bucket.Create(ctx, projectID, nil); err != nil {
		return err
	}

	fmt.Printf("Bucket: %v created.\n", bucket.BucketName())
	return nil
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
