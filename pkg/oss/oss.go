package oss

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"github.com/rs/xid"
)

// The above type defines an interface for interacting with a client that can perform operations on
// objects in a storage system.
type Clienter interface {
	// The `PutObject` function is used to upload an object to a storage system. It takes the following
	// parameters:
	PutObject(ctx context.Context, bucketName, objectName string, reader io.ReadSeeker, objectSize int64, options ...PutObjectOptions) (*PutObjectResult, error)
	// The `GetObject` function is used to retrieve an object from a storage system.
	GetObject(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	// The `RemoveObject` function is used to delete an object from a storage system.
	RemoveObject(ctx context.Context, bucketName, objectName string) error
	// The `StatObject` function is used to retrieve information about an object in a storage system.
	StatObject(ctx context.Context, bucketName, objectName string) (*ObjectStat, error)
}

// The `PutObjectOptions` type is used to specify options for putting an object, including content type
// and user metadata.
type PutObjectOptions struct {
	ContentType  string
	UserMetadata map[string]string
}

// The type PutObjectResult represents the result of a PUT request to upload an object.
type PutObjectResult struct {
	URL  string
	ETag string
}

// The ObjectStat type represents the metadata of an object, including its key, ETag, last modified
// time, size, content type, and user-defined metadata.
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
