package bucket

import (
	"context"
	"errors"
	"testing"

	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
)

func TestBucket_Create(t *testing.T) {
	t.Run("should create bucket", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		s3Mock := NewMockbucketHandler(ctrl)

		b := Bucket{
			s3: s3Mock,
		}

		output := &jam.OutputParam{}
		input := &jam.CreationParam{
			ID:     "ID_12345",
			Bucket: struct{ ID string }{"bucket_ID_12345"},
		}

		expectedCallParam := &s3.CreateBucketInput{
			ACL:    aws.String("private"),
			Bucket: aws.String("bucket_ID_12345"),
		}

		s3Mock.EXPECT().CreateBucket(expectedCallParam).Times(1)
		_, err := b.Create(context.Background(), input, output)
		if err != nil {
			t.Errorf("Bucket.Create() error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("should mutate output param", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		s3Mock := NewMockbucketHandler(ctrl)

		b := Bucket{
			s3: s3Mock,
		}

		output := &jam.OutputParam{}
		input := &jam.CreationParam{
			ID:     "ID_12345",
			Bucket: struct{ ID string }{"bucket_ID_12345"},
		}

		expectedCallParam := &s3.CreateBucketInput{
			ACL:    aws.String("private"),
			Bucket: aws.String("bucket_ID_12345"),
		}

		s3Mock.EXPECT().CreateBucket(expectedCallParam).Times(1)
		b.Create(context.Background(), input, output)
		if output.BucketID != "bucket_ID_12345" {
			t.Errorf("Bucket.Create() output.BucketID = %v, expected %v", output.BucketID, "bucket_ID_12345")
			return
		}
	})

	t.Run("should not create bucket and return err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		s3Mock := NewMockbucketHandler(ctrl)

		b := Bucket{
			s3: s3Mock,
		}

		output := &jam.OutputParam{}
		input := &jam.CreationParam{
			ID: "ID_12345",
		}

		expectedCallParam := &s3.CreateBucketInput{
			ACL:    aws.String("private"),
			Bucket: aws.String(""),
		}

		s3Mock.EXPECT().CreateBucket(expectedCallParam).Times(1).Return(nil, errors.New("Cannot create empty bucket"))
		_, err := b.Create(context.Background(), input, output)
		if err == nil {
			t.Errorf("Bucket.Create() there should be an error")
			return
		}
	})
}

func TestBucket_Destroy(t *testing.T) {
	t.Run("should delete bucket without emptying it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		s3Mock := NewMockbucketHandler(ctrl)

		b := Bucket{
			s3: s3Mock,
		}

		input := &jam.DeletionParam{
			BucketID: "bucket_ID_12345",
		}

		expectedCallDeleteParam := &s3.DeleteBucketInput{
			Bucket: aws.String("bucket_ID_12345"),
		}
		expectedCallListParam := &s3.ListObjectsInput{
			Bucket: aws.String("bucket_ID_12345"),
		}

		s3Mock.EXPECT().DeleteBucket(expectedCallDeleteParam).Times(1)
		s3Mock.EXPECT().ListObjects(expectedCallListParam).Times(1).Return(&s3.ListObjectsOutput{Contents: []*s3.Object{}}, nil)

		err := b.Destroy(context.Background(), input)
		if err != nil {
			t.Errorf("Bucket.Create() there should be no error: %s", err)
			return
		}
	})

	t.Run("should delete bucket after emptying", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		s3Mock := NewMockbucketHandler(ctrl)

		b := Bucket{
			s3: s3Mock,
		}

		input := &jam.DeletionParam{
			BucketID: "bucket_ID_12345",
		}

		expectedCallDeleteParam := &s3.DeleteBucketInput{
			Bucket: aws.String("bucket_ID_12345"),
		}
		expectedCallListParam := &s3.ListObjectsInput{
			Bucket: aws.String("bucket_ID_12345"),
		}

		del := []*s3.ObjectIdentifier{
			{
				Key: aws.String("test"),
			},
		}

		expectedCallDeleteObjParam := &s3.DeleteObjectsInput{
			Bucket: aws.String("bucket_ID_12345"),
			Delete: &s3.Delete{Objects: del},
		}

		returnListObject := &s3.ListObjectsOutput{
			IsTruncated: aws.Bool(false),
			Contents: []*s3.Object{
				{
					Key: aws.String("test"),
				},
			},
		}

		s3Mock.EXPECT().DeleteBucket(expectedCallDeleteParam).Times(1)
		s3Mock.EXPECT().ListObjects(expectedCallListParam).Times(1).Return(returnListObject, nil)
		s3Mock.EXPECT().DeleteObjects(expectedCallDeleteObjParam).Times(1)

		err := b.Destroy(context.Background(), input)
		if err != nil {
			t.Errorf("Bucket.Create() there should be no error: %s", err)
			return
		}
	})
}
