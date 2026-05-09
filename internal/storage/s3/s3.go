package s3

import (
	"context"
	"errors"
	"io"
	"time"

	"miso/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage struct {
	client         *s3.Client
	presignClient  *s3.PresignClient
	bucket         string
	requestTimeout time.Duration
	transferClient *transfermanager.Client
}

func New(config config.S3, sdkConfig aws.Config) *Storage {
	session := s3.NewFromConfig(sdkConfig)
	return &Storage{
		client:         session,
		presignClient:  s3.NewPresignClient(session),
		bucket:         config.Bucket,
		transferClient: transfermanager.New(session),
	}
}

func (s *Storage) requestContext() (context.Context, context.CancelFunc) {
	if s.requestTimeout > 0 {
		return context.WithTimeout(context.Background(), s.requestTimeout)
	}
	return context.Background(), func() {}
}

func (s *Storage) GetBuffer(key string) ([]byte, error) {
	var nsk *types.NoSuchKey

	if len(key) <= 0 {
		return nil, nil
	}
	ctx, cancel := s.requestContext()
	defer cancel()

	out, err := s.transferClient.GetObject(ctx, &transfermanager.GetObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(key),
	})
	if errors.As(err, &nsk) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if c, ok := out.Body.(io.Closer); ok {
		defer func() { _ = c.Close() }()
	}

	return io.ReadAll(out.Body)
}

func (s *Storage) GetStream(key string) (io.ReadCloser, error) {
	var nsk *types.NoSuchKey

	if len(key) <= 0 {
		return nil, nil
	}
	ctx, cancel := s.requestContext()
	defer cancel()

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(key),
	})
	if errors.As(err, &nsk) {
		return nil, nil
	}

	return resp.Body, err
}

func (s *Storage) Put(key string, data io.Reader) error {
	if len(key) <= 0 {
		return nil
	}

	ctx, cancel := s.requestContext()
	defer cancel()

	_, err := s.transferClient.UploadObject(ctx, &transfermanager.UploadObjectInput{
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
	if len(path) <= 0 {
		return nil, nil
	}

	ctx, cancel := s.requestContext()
	defer cancel()

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: aws.String(path),
	})

	var objects []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			objects = append(objects, aws.ToString(obj.Key))
		}
	}

	return objects, nil
}

func (s *Storage) GetPresignedURL(key string) (string, error) {
	if len(key) <= 0 {
		return "", nil
	}

	ctx, cancel := s.requestContext()
	defer cancel()

	presignResult, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(10 * time.Minute)
	})
	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}
