package filesystem

import (
	"errors"
	"io/fs"
	"kalisto/src/models"
	"os"
	"path/filepath"
	"strings"
)

type ProtoSearchResult struct {
	AbsoluteDirPath    string
	RelativeProtoPaths []string
}

// SearchProtoFiles function will find all .proto files by the given path.
func SearchProtoFiles(path string) (ProtoSearchResult, error) {
	result := ProtoSearchResult{}

	// Check if the path is absolute
	if !filepath.IsAbs(path) {
		return result, errors.New("path must be absolute")
	}

	// Check if the path is a directory or a file
	info, err := os.Stat(path)
	if err != nil {
		return result, err
	}

	// This is a file
	if !info.IsDir() {
		if !strings.HasSuffix(info.Name(), ".proto") {
			return result, errors.New("chosen file is not a proto file")
		}

		result.AbsoluteDirPath = filepath.Dir(path)
		result.RelativeProtoPaths = []string{filepath.Base(path)}

		return result, nil
	}

	// This is a directory, find all .proto files recursively
	result.AbsoluteDirPath = path
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".proto") {
			rel, err := filepath.Rel(result.AbsoluteDirPath, path)
			if err != nil {
				return err
			}
			result.RelativeProtoPaths = append(result.RelativeProtoPaths, rel)
		}
		return nil
	})

	if err != nil {
		return result, err
	}

	if len(result.RelativeProtoPaths) == 0 {
		return result, models.ErrNoProtoFilesFound
	}

	return result, nil
}
