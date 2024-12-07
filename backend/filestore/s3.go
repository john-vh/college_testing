package filestore

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Store struct {
	logger   *slog.Logger
	s3Client *s3.S3
	bucket   *string
	region   *string
}

func NewS3ImageStore(profile string, bucket string, region string) (*S3Store, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile: profile,
	})
	if err != nil {
		return nil, err
	}
	return &S3Store{
		bucket:   aws.String(bucket),
		region:   aws.String(region),
		s3Client: s3.New(s),
	}, nil
}

func (s *S3Store) UploadObject(key string, f io.ReadSeeker) (err error) {
	_, err = s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *S3Store) DeleteObject(key string) (err error) {
	_, err = s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
	})
	return err
}

func (s *S3Store) GetURI(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", *s.bucket, *s.region, key)
}

func (s *S3Store) GetKey(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]

}
