package oss

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3ClientConfig struct {
	Domain          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Prefix          string
}

var _ IClient = (*S3Client)(nil)

type S3Client struct {
	config  S3ClientConfig
	session *session.Session
	client  *s3.S3
}

func NewS3Client(config S3ClientConfig) (*S3Client, error) {
	awsConfig := aws.NewConfig()
	awsConfig.WithRegion(config.Region)
	awsConfig.WithCredentials(credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""))
	session, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		config:  config,
		session: session,
		client:  s3.New(session),
	}, nil
}

func (c *S3Client) PutObject(ctx context.Context, bucketName, objectName string, reader io.ReadSeeker, objectSize int64, options ...PutObjectOptions) (*PutObjectResult, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	var opt PutObjectOptions
	if len(options) > 0 {
		opt = options[0]
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	input := &s3.PutObjectInput{
		Bucket:             aws.String(bucketName),
		Key:                aws.String(objectName),
		Body:               reader,
		ContentType:        aws.String(opt.ContentType),
		ContentDisposition: aws.String("inline"),
		ACL:                aws.String("public-read"),
	}

	if len(opt.UserMetadata) > 0 {
		input.Metadata = make(map[string]*string)
		for k, v := range opt.UserMetadata {
			input.Metadata[k] = aws.String(v)
		}
	}

	output, err := c.client.PutObject(input)
	if err != nil {
		return nil, err
	}

	return &PutObjectResult{
		URL:  c.config.Domain + "/" + objectName,
		Key:  *input.Key,
		ETag: *output.ETag,
		Size: objectSize,
	}, nil
}

func (c *S3Client) GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	output, err := c.client.GetObject(input)
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

func (c *S3Client) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	_, err := c.client.DeleteObject(input)
	return err
}

func (c *S3Client) RemoveObjectByURL(ctx context.Context, urlStr string) error {
	prefix := c.config.Domain + "/"
	if !strings.HasPrefix(urlStr, prefix) {
		return nil
	}

	objectName := strings.TrimPrefix(urlStr, prefix)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(c.config.BucketName),
		Key:    aws.String(objectName),
	}

	_, err := c.client.DeleteObject(input)
	return err
}

func (c *S3Client) StatObjectByURL(ctx context.Context, urlStr string) (*ObjectStat, error) {
	prefix := c.config.Domain + "/"
	if !strings.HasPrefix(urlStr, prefix) {
		return nil, nil
	}

	objectName := strings.TrimPrefix(urlStr, prefix)
	return c.StatObject(ctx, c.config.BucketName, objectName)
}

func (c *S3Client) StatObject(ctx context.Context, bucketName, objectName string) (*ObjectStat, error) {
	if bucketName == "" {
		bucketName = c.config.BucketName
	}

	objectName = formatObjectName(c.config.Prefix, objectName)
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	output, err := c.client.HeadObject(input)
	if err != nil {
		return nil, err
	}

	var metadata map[string]string
	if output.Metadata != nil {
		metadata = make(map[string]string)
		for k, v := range output.Metadata {
			metadata[k] = *v
		}
	}

	return &ObjectStat{
		Key:          objectName,
		ETag:         *output.ETag,
		LastModified: *output.LastModified,
		Size:         *output.ContentLength,
		ContentType:  *output.ContentType,
		UserMetadata: metadata,
	}, nil
}
