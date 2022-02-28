package pa

import (
	"context"
	"io"
)

// FileService represents a service which manages files in the system.
// Should be used to create / delete files in a served fs.
type FileService interface {
	// CreateFile creates a new file.
	CreateFile(ctx context.Context, path string, content io.Reader) error

	// DeletePath deletes path.
	DeleteFile(ctx context.Context, path string) error
}
