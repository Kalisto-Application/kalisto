package filesystem

import (
	"kalisto/src/models"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchProtoFiles(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		want    ProtoSearchResult
		errorIs error
		wantErr error
	}{
		{name: "file or directory not found", path: path.Join(wd, "fake_path"), want: ProtoSearchResult{}, errorIs: os.ErrNotExist},
		{name: "path must be absolute", path: "testdata/search_proto", want: ProtoSearchResult{}, wantErr: &models.ErrorFileMustBeAbsolute{File: "testdata/search_proto"}},
		{name: "no proto files in directory", path: path.Join(wd, "testdata/search_proto/empty_dir"), errorIs: models.ErrNoProtoFilesFound},
		{name: "happy path: search in directory", path: path.Join(wd, "testdata/search_proto"), want: ProtoSearchResult{
			AbsoluteDirsPath:   []string{path.Join(wd, "testdata/search_proto")},
			RelativeProtoPaths: []string{"file1.proto", "sub_dir/file2.proto", "sub_dir/sub_dir2/file3.proto"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchProtoFiles([]string{tt.path})

			if tt.errorIs != nil {
				require.ErrorIs(t, err, tt.errorIs)
			} else if tt.wantErr != nil {
				require.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.AbsoluteDirsPath, got.AbsoluteDirsPath)
				require.ElementsMatch(t, tt.want.RelativeProtoPaths, got.RelativeProtoPaths)
			}
		})
	}
}
