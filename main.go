package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jankaczmarski/go-s3/internal/aws"
	"github.com/jankaczmarski/go-s3/internal/gcp"
	mytypes "github.com/jankaczmarski/go-s3/pkg/types"

	"github.com/jankaczmarski/go-s3/pkg/types"
)

type CloudProvider string

const (
	projectID                = "go-s3-play"
	awsRegion                = "eu-central-1"
	awsProfile               = "jank-private"
	AWS        CloudProvider = "AWS"
	GCP        CloudProvider = "GCP"
)

func run(ctx context.Context, worker types.StorageWorker) {
	newBucketName := "test-switch-case-gos3"

	bucket, err := worker.CreateBucket(ctx, newBucketName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bucket Created: %s\n", bucket.Name)

	buckets, err := worker.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("Failed to listBuckets: %v", err)
	}

	for i, tmpBucketName := range buckets {
		fmt.Printf("ID-%d: BucketName: %s\n", i, tmpBucketName)
	}
}

// User gets github.com/jankaczmarski/go-s3/xs3 lib and he want to ListBuckets on provider
// He sets provider on which he want to run in switch case and runs ListBuckets(ctx, worker) -> he gets buckets he want
// The same with uploading data and receiving -> Is this approach correct nad what library do we want to expose? just the
// Maybe put each provider in one file in pkg folder, so we have package xs3 and aws.go and gcp.go in this folder? But what do we need internal for then, maybe it's not needed so much?
func main() {
	// Define a flag for cloudProvider
	cloudProviderFlag := flag.String("cloudProvider", "aws", "The cloud provider to use (aws or azure)")
	flag.Parse()

	// Use the flag value to set the cloudProvider variable
	cloudProvider := CloudProvider(*cloudProviderFlag)
	var err error
	var worker mytypes.StorageWorker
	ctx := context.Background()

	switch cloudProvider {
	case AWS:
		worker, err = aws.NewWorker(ctx, awsProfile)
		if err != nil {
			log.Fatal(err)
		}
	case GCP:
		worker, err = gcp.NewWorker(ctx, projectID)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("No valid provider passed")
	}
	run(ctx, worker)
}
