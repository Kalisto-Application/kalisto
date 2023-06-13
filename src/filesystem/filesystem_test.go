package filesystem

import (
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
		wantErr string
	}{
		{name: "file or directory not found", path: path.Join(wd, "fake_path"), want: ProtoSearchResult{}, wantErr: "no such file or directory"},
		{name: "path must be absolute", path: "testdata/search_proto", want: ProtoSearchResult{}, wantErr: "path must be absolute"},
		{name: "no proto files in directory", path: path.Join(wd, "testdata/search_proto/empty_dir"), wantErr: "no proto files found"},
		{name: "single file is not a proto file", path: path.Join(wd, "testdata/search_proto/not_proto_file.txt"), wantErr: "not a proto file"},
		{name: "happy path: single file", path: path.Join(wd, "testdata/search_proto/file1.proto"), want: ProtoSearchResult{
			AbsoluteDirPath:    path.Join(wd, "testdata/search_proto"),
			RelativeProtoPaths: []string{"file1.proto"},
		}},
		{name: "happy path: search in directory", path: path.Join(wd, "testdata/search_proto"), want: ProtoSearchResult{
			AbsoluteDirPath:    path.Join(wd, "testdata/search_proto"),
			RelativeProtoPaths: []string{"file1.proto", "sub_dir/file2.proto", "sub_dir/sub_dir2/file3.proto"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchProtoFiles(tt.path)
			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.AbsoluteDirPath, got.AbsoluteDirPath)
			require.ElementsMatch(t, tt.want.RelativeProtoPaths, got.RelativeProtoPaths)
		})
	}
}
