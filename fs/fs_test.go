package fs_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/fs"
)

// CreateWitCleanup will create root and then delete the root path and all its children.
func CreateWithCleanup(root string) (pa.FileService, func()) {
	return fs.NewFileService(root), func() { os.RemoveAll(root) }
}

func TestNewFileService(t *testing.T) {
	_, cln := CreateWithCleanup("./foo/")
	defer cln()

	if _, err := os.Stat("./foo/"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateFile(t *testing.T) {
	fs, cln := CreateWithCleanup("./foo")
	defer cln()

	adminUsrCtx := pa.NewContextWithUser(context.Background(), &pa.User{IsAdmin: true})

	// create file.
	if err := fs.CreateFile(adminUsrCtx, "/bar/file.txt", strings.NewReader("FOO BAR")); err != nil {
		t.Fatal(err)
	} else if err := fs.CreateFile(adminUsrCtx, "/file.txt", strings.NewReader("FOO BAR")); err != nil {
		t.Fatal(err)
	}

	// assert creation.
	if _, err := os.Stat("./foo/bar/file.txt"); err != nil {
		t.Fatal(err)
	} else if _, err := os.Stat("./foo/file.txt"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteFile(t *testing.T) {
	fs, _ := CreateWithCleanup("./foo")

	adminUsrCtx := pa.NewContextWithUser(context.Background(), &pa.User{IsAdmin: true})

	// create file.
	MustCreateFile(t, fs, adminUsrCtx, "/baz/file.txt", strings.NewReader("FOOOOOOO"))

	t.Run("Ok Delete Call", func(t *testing.T) {
		if err := fs.DeleteFile(adminUsrCtx, "/baz/file.txt"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Bad Delete Call (Not Found)", func(t *testing.T) {
		if err := fs.DeleteFile(adminUsrCtx, "/baz/file.txt"); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func MustCreateFile(t *testing.T, fs pa.FileService, ctx context.Context, path string, content io.Reader) {
	t.Helper()
	if err := fs.CreateFile(ctx, path, content); err != nil {
		t.Fatal(err)
	}
}
