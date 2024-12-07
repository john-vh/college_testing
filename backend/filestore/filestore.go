package filestore

import "io"

type FileStore interface {
	UploadObject(key string, f io.ReadSeeker) error
}
