// Sample storage-quickstart creates a Google Cloud Storage bucket.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	// context for working with API's, sender sends context and receiver receieves it,
	// imo it's a context of certain api calls, more info here: https://pkg.go.dev/context
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "go-s3-play"

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name for the new bucket.
	bucketName := "jank-go-s3-bucket"
	bucket := client.Bucket(bucketName)

	buckets, err := listBuckets(client, ctx, projectID)

	if !slices.Contains(buckets, bucketName) {
		err := createBucket(bucket, ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}

	err = uploadFile(log.Writer(), bucketName, "example-data/notes.txt")
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
}

func createBucket(bucket *storage.BucketHandle, ctx context.Context, projectID string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		return err
	}

	fmt.Printf("Bucket: %v created.\n", bucket.BucketName())
	return nil
}

func listBuckets(client *storage.Client, ctx context.Context, projectID string) ([]string, error) {
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

// uploadFile uploads an object.
// NOTE: For now this function only can create new objects, it can't modify existing ones
func uploadFile(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Open local file.
	f, err := os.Open("notes.txt")
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return nil
}
