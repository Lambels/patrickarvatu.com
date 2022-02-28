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

func NewImageService(root string) *FileService {
	os.MkdirAll(root, 0755)

	return &FileService{
		root: root,
	}
}

func (s *FileService) CreateFile(ctx context.Context, path string, content io.Reader) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	// cut trailing "/" and splice the path on "/".
	splicedPath := strings.Split(strings.TrimLeft(string(filepath.Separator), path), string(filepath.Separator))

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

func (s *FileService) DeleteFile(ctx context.Context, path string) error {
	if !pa.IsAdminContext(ctx) {
		return pa.Errorf(pa.EUNAUTHORIZED, "user isnt admin.")
	}

	// remove file.
	if err := os.Remove(path); err != nil {
		// parse error.
		if os.IsNotExist(err) {
			return pa.Errorf(pa.ENOTFOUND, "file doesent exist")
		}
		return err
	}
	return nil
}

// TODO: Implement.
func pathExists(path string) bool {

}
