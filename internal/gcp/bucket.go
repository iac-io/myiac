package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
)

const tfstate = "-tfstate-"

func CreateGCSBucket(projectID string, e string) error {
	// Setup context, client and bucket name
	bucketName := projectID + tfstate + e
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Setup client bucket to work from
	bucket := client.Bucket(bucketName)

	buckets := client.Buckets(ctx, projectID)
	for {
		if bucketName == "" {
			return fmt.Errorf("BucketName entered is empty %v.", bucketName)
		}
		attrs, err := buckets.Next()
		// Assume bucket not found if at Iterator end and create
		if err == iterator.Done {
			// Create bucket
			if err := bucket.Create(ctx, projectID, &storage.BucketAttrs{
				Location: "EU",
			}); err != nil {
				return fmt.Errorf("Failed to create bucket: %v", err)
			}
			log.Printf("Bucket %v created.\n", bucketName)
			return nil
		}
		if err != nil {
			return fmt.Errorf("Issues setting up Bucket(%q).Objects(): %v. Double check project id.", attrs.Name, err)
		}
		if attrs.Name == bucketName {
			log.Printf("Bucket %v exists.\n", bucketName)
			return nil
		}
	}
}

func DeleteGCSBucket(projectID string, e string) error {
	// Setup context, client and bucket name
	bucketName := projectID + tfstate + e
	fmt.Printf("Deleting Bucket: %v\n", bucketName)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("GCP: Could not connect to Google Cloud. Error: %v", err)
	}
	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return fmt.Errorf("GCP: Stroage: Could not delete Bucket: %v.Error: %v", bucketName, err)
	}
	fmt.Printf("Bucket %v deleted!\n", bucketName)
	return nil
}