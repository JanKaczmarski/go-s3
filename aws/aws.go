package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// TODO: Make this Client field unexported and create Close method on AwsWorker type
type AwsWorker struct {
	Client *s3.Client
}

func NewWorker(ctx context.Context) (*AwsWorker, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("jank-private"),
	)
	if err != nil {
		return nil, err
	}
	return &AwsWorker{
		Client: s3.NewFromConfig(cfg),
	}, nil
}

func (worker *AwsWorker) ListBuckets(ctx context.Context) ([]string, error) {
	var buckets []string
	bucketPaginator := s3.NewListBucketsPaginator(worker.Client, &s3.ListBucketsInput{})

	for bucketPaginator.HasMorePages() {
		output, err := bucketPaginator.NextPage(ctx)
		if err != nil {
			log.Printf("Couldn't list buckets for your account, because: %v\n", err)
			return nil, err
		}
		for _, bucket := range output.Buckets {
			buckets = append(buckets, *bucket.Name)
		}
	}
	return buckets, nil
}
