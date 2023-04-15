package oss

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/rs/xid"
)

// Object storage client interface
type Clienter interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.ReadSeeker, objectSize int64, options ...PutObjectOptions) (*PutObjectResult, error)
	GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	RemoveObject(ctx context.Context, bucketName, objectName string) error
	StatObject(ctx context.Context, bucketName, objectName string) (*ObjectStat, error)
}

type PutObjectOptions struct {
	ContentType  string
	UserMetadata map[string]string
}

type PutObjectResult struct {
	URL  string
	ETag string
}

type ObjectStat struct {
	Key          string
	ETag         string
	LastModified time.Time
	Size         int64
	ContentType  string
	UserMetadata map[string]string
}

func (a *ObjectStat) GetName() string {
	if name, ok := a.UserMetadata["name"]; ok {
		return name
	}
	return filepath.Base(a.Key)
}

func formatObjectName(objectName string) string {
	if objectName == "" {
		return xid.New().String()
	}
	if objectName[0] == '/' {
		objectName = objectName[1:]
	}
	return objectName
}
