package s3

import (
	"context"
	"errors"
	"io"
	"time"

	"miso/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	client         *s3.Client
	bucket         string
	requestTimeout time.Duration
	downloader     *manager.Downloader
	uploader       *manager.Uploader
}

func New(config config.S3, sdkConfig aws.Config) *Storage {
	session := s3.NewFromConfig(sdkConfig)
	return &Storage{
		client:     session,
		bucket:     config.Bucket,
		downloader: manager.NewDownloader(session),
		uploader:   manager.NewUploader(session),
	}
}

func (s *Storage) requestContext() (context.Context, context.CancelFunc) {
	if s.requestTimeout > 0 {
		return context.WithTimeout(context.Background(), s.requestTimeout)
	}
	return context.Background(), func() {}
}

func (s *Storage) Get(key string) ([]byte, error) {
	var nsk *types.NoSuchKey

	if len(key) <= 0 {
		return nil, nil
	}
	ctx, cancel := s.requestContext()
	defer cancel()

	buf := manager.NewWriteAtBuffer([]byte{})

	_, err := s.downloader.Download(ctx, buf, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(key),
	})
	if errors.As(err, &nsk) {
		return nil, nil
	}

	return buf.Bytes(), err
}

func (s *Storage) Put(key string, data io.Reader) error {
	if len(key) <= 0 {
		return nil
	}

	ctx, cancel := s.requestContext()
	defer cancel()

	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(key),
		Body:   data,
	})

	return err
}

func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}

	ctx, cancel := s.requestContext()
	defer cancel()

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(key),
	})

	return err
}

func (s *Storage) List(path string) ([]string, error) {
	var objects []string

	if len(path) <= 0 {
		return nil, nil
	}

	ctx, cancel := s.requestContext()
	defer cancel()

	listObjectsInput := s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: aws.String(path),
	}

	resp, err := s.client.ListObjectsV2(ctx, &listObjectsInput)
	if err != nil {
		return nil, err
	}

	for _, object := range resp.Contents {
		objects = append(objects, *object.Key)
	}

	return objects, err
}
