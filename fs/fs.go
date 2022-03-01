package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	pa "github.com/Lambels/patrickarvatu.com"
)

var _ pa.FileService = (*FileService)(nil)

type FileService struct {
	root string
}

// NewFileService returns a new FileService with root path set to root.
func NewFileService(root string) *FileService {
	os.MkdirAll(root, 0755)

	return &FileService{
		root: root,
	}
}

// CreateFile creates a new path ending with the file ie: "/bar/baz/file.txt" will create
// a file bar with children baz if they dont exist and a file.txt in baz.
func (s *FileService) CreateFile(ctx context.Context, path string, content io.Reader) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	// cut trailing "/" and splice the path on "/".
	splicedPath := strings.Split(strings.TrimLeft(path, string(filepath.Separator)), string(filepath.Separator))

	// we want to ensure the path leading to the file exists before opening it.
	if len(splicedPath) > 1 {
		// cut the file name from the path.
		dirPath := strings.Join(splicedPath[:len(splicedPath)-1], "")

		// path doesent exist: create.
		if !pathExists(filepath.Join(s.root, dirPath)) {
			// create dir and dir parents.
			if err := os.MkdirAll(filepath.Join(s.root, dirPath), 0777); err != nil {
				return err
			}
		}
	}

	// create file with already existing path.
	f, err := os.Create(filepath.Join(s.root, path))
	if err != nil {
		return err
	}
	defer f.Close()

	// copy content into file.
	_, err = io.Copy(f, content)
	return err
}

// DeleteFile deletes the file path in the file system.
// returns ENOTFOUND if file isnt found.
func (s *FileService) DeleteFile(ctx context.Context, path string) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	// remove file.
	if err := os.Remove(filepath.Join(s.root, path)); err != nil {
		// parse error.
		if os.IsNotExist(err) {
			return pa.Errorf(pa.ENOTFOUND, "file doesent exist.")
		}
		return err
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
