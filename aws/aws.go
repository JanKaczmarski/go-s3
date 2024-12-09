package aws

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jankaczmarski/go-s3/helpers"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type BucketBasics struct {
	S3Client *s3.Client
}

func aws() {
	ctx := context.Background()
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("jank-private"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	basics := BucketBasics{s3.NewFromConfig(cfg)}

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := basics.ListBuckets(ctx)
	if err != nil {
		log.Fatalf("ListBuckets resulted in error: %v", err)
	}

	for i, bucket := range output {
		bucketRegion := helpers.DefaultString(bucket.BucketRegion, "Unaccessable")
		fmt.Printf("ID-%v: Bucket Name: %s, Bucket Region: %s\n", i, *bucket.Name, bucketRegion)
	}
}

func (basics BucketBasics) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	var err error
	var output *s3.ListBucketsOutput
	var buckets []types.Bucket
	bucketPaginator := s3.NewListBucketsPaginator(basics.S3Client, &s3.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				fmt.Println("You don't have permission to list buckets for this account.")
				err = apiErr
			} else {
				log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
			}
			break
		} else {
			buckets = append(buckets, output.Buckets...)
		}
	}
	return buckets, err
}
