package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	// partSize bounds memory per in-flight part; with 10000 parts max it allows dumps up to ~160 GB.
	partSize   = 16 << 20
	awsS3Host  = "s3.amazonaws.com"
	uploadType = "application/gzip"
)

type s3Store struct {
	client *minio.Client
	bucket string
}

func newS3Store(c *config) (*s3Store, error) {
	host, secure := awsS3Host, true

	if c.s3Endpoint != "" {
		raw := c.s3Endpoint
		if !strings.Contains(raw, "://") {
			raw = "https://" + raw
		}

		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			return nil, fmt.Errorf("S3_ENDPOINT %q is not a valid URL", c.s3Endpoint)
		}

		host, secure = u.Host, u.Scheme == "https"
	}

	lookup := minio.BucketLookupAuto
	if c.s3PathStyle {
		lookup = minio.BucketLookupPath
	}

	client, err := minio.New(host, &minio.Options{
		Creds:        credentials.NewStaticV4(c.s3AccessKey, c.s3SecretKey, ""),
		Secure:       secure,
		Region:       c.s3Region,
		BucketLookup: lookup,
	})
	if err != nil {
		return nil, fmt.Errorf("s3 client: %w", err)
	}

	return &s3Store{client: client, bucket: c.s3Bucket}, nil
}

// Put streams r as a multipart upload; an unknown size (-1) is required as the dump length is not known upfront.
func (s *s3Store) Put(ctx context.Context, key string, r io.Reader) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, r, -1, minio.PutObjectOptions{
		ContentType: uploadType,
		PartSize:    partSize,
		NumThreads:  1,
	})

	return err //nolint:wrapcheck // wrapped by caller
}

func (s *s3Store) Stat(ctx context.Context, key string) (int64, error) {
	info, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return 0, err //nolint:wrapcheck // wrapped by caller
	}

	return info.Size, nil
}

func (s *s3Store) List(ctx context.Context, prefix string) ([]object, error) {
	var out []object

	for o := range s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if o.Err != nil {
			return nil, o.Err
		}

		out = append(out, object{key: o.Key, size: o.Size, modified: o.LastModified})
	}

	return out, nil
}

func (s *s3Store) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err //nolint:wrapcheck // wrapped by caller
	}

	if _, err = obj.Stat(); err != nil {
		_ = obj.Close()

		return nil, err //nolint:wrapcheck // wrapped by caller
	}

	return obj, nil
}

func (s *s3Store) Delete(ctx context.Context, keys []string) error {
	in := make(chan minio.ObjectInfo, len(keys))
	for _, k := range keys {
		in <- minio.ObjectInfo{Key: k}
	}

	close(in)

	var errs []error

	for e := range s.client.RemoveObjects(ctx, s.bucket, in, minio.RemoveObjectsOptions{}) {
		errs = append(errs, fmt.Errorf("%s: %w", e.ObjectName, e.Err))
	}

	return errors.Join(errs...)
}
