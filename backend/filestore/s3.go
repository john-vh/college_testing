package filestore

import (
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Store struct {
	logger   *slog.Logger
	s3Client *s3.S3
	bucket   *string
}

func NewS3ImageStore(bucket string, region string) (*S3Store, error) {
	s, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &S3Store{
		bucket:   aws.String(bucket),
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
