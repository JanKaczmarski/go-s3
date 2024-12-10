package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jankaczmarski/go-s3/internal/aws"
	"github.com/jankaczmarski/go-s3/pkg/types"
)

const (
	projectID = "go-s3-play"
	awsRegion = "eu-central-1"
)

func run(ctx context.Context, worker types.StorageWorker) {
	buckets, err := worker.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("Failed to listBuckets: %v", err)
	}

	for i, tmpBucketName := range buckets {
		fmt.Printf("ID-%d: BucketName: %s\n", i, tmpBucketName)
	}
}

func main() {
	// GCP testing
	ctx := context.Background()
	//worker, err := gcp.NewWorker(ctx, projectID)
	//if err != nil {
	//	log.Fatalf("Failed to create NewWorker for gcp: %v", err)
	//}

	//defer worker.Close()

	//fmt.Println("Running GCP worker")
	//run(ctx, worker)
	//fmt.Println("GCP worker finished running")

	// AWS testing
	fmt.Println("Setup AWS worker")

	ctxAws := context.Background()
	workerAws, err := aws.NewWorker(ctxAws)
	if err != nil {
		log.Fatalf("Failed to create NewWorker for AWS: %v", err)
	}

	bucketNameAws := "testing-bucket-go-s3-01"
	log.Printf("Creating aws bucket: %s\n", bucketNameAws)

	_, err = aws.NewAwsBucket(ctx, bucketNameAws, awsRegion, workerAws)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created bucket: %s\n", bucketNameAws)

	fmt.Println("Running AWS worker")
	run(ctxAws, workerAws)
	fmt.Println("AWS worker finished running")
}
