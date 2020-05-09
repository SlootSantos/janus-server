package bucket

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (b *Bucket) emptyBucket(bucketID string) error {
	log.Println("removing objects from S3 bucket : ", bucketID)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketID),
	}
	for {
		//Requesting for batch of objects from s3 bucket
		objects, err := b.s3.ListObjects(params)
		if err != nil {
			return err
		}
		//Checks if the bucket is already empty
		if len((*objects).Contents) == 0 {
			log.Println("Bucket is already empty")
			return nil
		}
		log.Println("First object in batch | ", *(objects.Contents[0].Key))

		//creating an array of pointers of ObjectIdentifier
		objectsToDelete := make([]*s3.ObjectIdentifier, 0, 1000)
		for _, object := range (*objects).Contents {
			obj := s3.ObjectIdentifier{
				Key: object.Key,
			}
			objectsToDelete = append(objectsToDelete, &obj)
		}
		//Creating JSON payload for bulk delete
		deleteArray := s3.Delete{Objects: objectsToDelete}
		deleteParams := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketID),
			Delete: &deleteArray,
		}
		//Running the Bulk delete job (limit 1000)
		_, err = b.s3.DeleteObjects(deleteParams)
		if err != nil {
			return err
		}
		if *(*objects).IsTruncated { //if there are more objects in the bucket, IsTruncated = true
			params.Marker = (*deleteParams).Delete.Objects[len((*deleteParams).Delete.Objects)-1].Key
			log.Println("Requesting next batch | ", *(params.Marker))
		} else { //if all objects in the bucket have been cleaned up.
			break
		}
	}
	log.Println("Emptied S3 bucket : ", bucketID)
	return nil
}
