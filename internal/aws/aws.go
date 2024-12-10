package aws

import (
	"context"
	"errors"
	"log"
	"time"

	mytypes "github.com/jankaczmarski/go-s3/pkg/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AwsWorker struct {
	client *s3.Client
}

type AwsBucket struct {
	bucketName string
}

func NewWorker(ctx context.Context) (*AwsWorker, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile("jank-private"),
	)
	if err != nil {
		return nil, err
	}
	return &AwsWorker{
		client: s3.NewFromConfig(cfg),
	}, nil
}

// TODO: jk: this func creates bucket if it's not present. Make distinct function for each functionality
// One for getting new bucket and one for creating bucket
func NewAwsBucket(ctx context.Context, bucketName, region string, worker *AwsWorker) (*AwsBucket, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := worker.client.CreateBucket(ctx, &s3.CreateBucketInput{
		// aws only accepts *string and this is how you do this per aws docs
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})

	// This is aws error handling style taken from docs: https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s.\n", bucketName)
			err = owned
		} else if errors.As(err, &exists) {
			log.Printf("Bucket %s already exists.\n", bucketName)
			err = exists
		}
	} else {
		err = s3.NewBucketExistsWaiter(worker.client).Wait(
			ctx, &s3.HeadBucketInput{Bucket: aws.String(bucketName)}, time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for bucket %s to exist.\n", bucketName)
		}
	}

	return &AwsBucket{
		bucketName: bucketName,
	}, err
}

func (worker *AwsWorker) CreateBucket(ctx context.Context, name, region string) (*mytypes.Bucket, error) {
	awsBucket, err := NewAwsBucket(ctx, name, region, worker)
	if err != nil {
		return nil, err
	}
	return &mytypes.Bucket{
		Name:   awsBucket.bucketName,
		Region: region,
	}, nil
}

func (worker *AwsWorker) ListBuckets(ctx context.Context) ([]string, error) {
	var buckets []string
	bucketPaginator := s3.NewListBucketsPaginator(worker.client, &s3.ListBucketsInput{})

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
