package filestore

import "io"

type FileStore interface {
	UploadObject(key string, f io.ReadSeeker) error
	DeleteObject(key string) error
	GetURI(key string) (URI string)
	GetKey(url string) (key string)
}
