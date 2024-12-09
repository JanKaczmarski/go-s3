package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jankaczmarski/go-s3/gcp"
)

const (
	projectID = "go-s3-play"
)

type storageWorker interface {
	ListBuckets(context.Context) ([]string, error)
}

func run(ctx context.Context, worker storageWorker) {
	buckets, err := worker.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("Failed to listBuckets: %v", err)
	}

	for i, tmpBucketName := range buckets {
		fmt.Printf("ID-%d: BucketName: %s\n", i, tmpBucketName)
	}
}

func main() {
	ctx := context.Background()
	worker, err := gcp.NewWorker(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create NewWorker for gcp: %v", err)
	}

	client := worker.Client
	defer client.Close()

	if err != nil {
		log.Fatalf("Creation of new GcpWorker failed with error: %v", err)
	}

	fmt.Println("Running worker")
	run(ctx, worker)
	fmt.Println("Worker finished running")
}
