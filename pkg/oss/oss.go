package oss

import (
	"context"
	"io"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/xid"
)

var (
	Ins  IClient
	once sync.Once
)

// Set the global oss client
func SetGlobal(h func() IClient) {
	once.Do(func() {
		Ins = h()
	})
}

// IClient is an interface for oss client
type IClient interface {
	PutObject(ctx context.Context, bucketName, objectName string, reader io.ReadSeeker, objectSize int64, options ...PutObjectOptions) (*PutObjectResult, error)
	GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	RemoveObject(ctx context.Context, bucketName, objectName string) error
	RemoveObjectByURL(ctx context.Context, urlStr string) error
	StatObject(ctx context.Context, bucketName, objectName string) (*ObjectStat, error)
	StatObjectByURL(ctx context.Context, urlStr string) (*ObjectStat, error)
}

// PutObjectOptions represents options specified by user for PutObject call
type PutObjectOptions struct {
	ContentType  string
	UserMetadata map[string]string
}

type PutObjectResult struct {
	URL  string `json:"url,omitempty"`
	Key  string `json:"key,omitempty"`
	ETag string `json:"e_tag,omitempty"`
	Size int64  `json:"size,omitempty"`
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

func formatObjectName(prefix, objectName string) string {
	if objectName == "" {
		objectName = xid.New().String()
	}
	if objectName[0] == '/' {
		objectName = objectName[1:]
	}
	if prefix != "" {
		objectName = prefix + "/" + objectName
	}
	return objectName
}
