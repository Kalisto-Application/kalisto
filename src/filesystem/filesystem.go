package filesystem

import (
	"io/fs"
	"kalisto/src/models"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type bufWork struct {
	Directories []string `yaml:"directories"`
}

type ProtoSearchResult struct {
	AbsoluteDirsPath   []string
	RelativeProtoPaths []string
	BufDirs            []string
}

// SearchProtoFiles function will find all .proto files by the given path.
func SearchProtoFiles(dirs []string) (ProtoSearchResult, error) {
	result := ProtoSearchResult{}

	for _, dir := range dirs {
		// Check if the path is absolute
		if !filepath.IsAbs(dir) {
			return result, &models.ErrorFileMustBeAbsolute{File: dir}
		}

		// Check if the path is a directory or a file
		// info, err := os.Stat(dir)
		// if err != nil {
		// 	return result, err
		// }
		result.BufDirs = append(result.BufDirs, readBufWorkDirs(dir)...)

		// // This is a file
		// if !info.IsDir() {
		// 	if !strings.HasSuffix(info.Name(), ".proto") {
		// 		return result, errors.New("chosen file is not a proto file")r
		// 	}

		// 	result.AbsoluteDirsPath = append(result.AbsoluteDirsPath, filepath.Dir(dir))
		// 	result.RelativeProtoPaths = append(result.RelativeProtoPaths, filepath.Base(dir))
		// 	continue
		// }

		// This is a directory, find all .proto files recursively
		result.AbsoluteDirsPath = append(result.AbsoluteDirsPath, dir)
		if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(d.Name(), ".proto") {
				rel, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}
				result.RelativeProtoPaths = append(result.RelativeProtoPaths, rel)
			}
			return nil
		}); err != nil {
			return result, err
		}
	}

	if len(result.RelativeProtoPaths) == 0 {
		return result, models.ErrNoProtoFilesFound
	}

	return result, nil
}

func readBufWorkDirs(path string) []string {
	f, err := os.Open(filepath.Join(path, "buf.work.yaml"))
	if err != nil {
		f, err = os.Open(filepath.Join(path, "buf.work.yml"))
		if err != nil {
			return nil
		}
	}

	var buf bufWork
	if err := yaml.NewDecoder(f).Decode(&buf); err != nil {
		return nil
	}

	return buf.Directories
}

func SearchOpenapiFiles(files []string) ([]string, error) {
	for _, filename := range files {
		if !path.IsAbs(filename) {
			return nil, &models.ErrorFileMustBeAbsolute{File: filename}
		}

		stat, err := os.Stat(filename)
		if err != nil {
			return nil, err
		}

		if stat.IsDir() {
			return nil, &models.ErrorOpenapiFileCantBeDir{File: filename}
		}
	}

	return files, nil
}
