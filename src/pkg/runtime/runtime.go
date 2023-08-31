package runtime

import (
	"context"

	rpkg "github.com/wailsapp/wails/v2/pkg/runtime"
)

type runtime struct {
}

func New() *runtime {
	return &runtime{}
}

func (r *runtime) MessageDialog(ctx context.Context, opts rpkg.MessageDialogOptions) (string, error) {
	if ctx == nil {
		return "", nil
	}
	return rpkg.MessageDialog(ctx, opts)
}

func (r *runtime) OpenDirectoryDialog(ctx context.Context, opts rpkg.OpenDialogOptions) (string, error) {
	if ctx == nil {
		return "", nil
	}
	return rpkg.OpenDirectoryDialog(ctx, opts)
}

type Runtime interface {
	MessageDialog(ctx context.Context, opts rpkg.MessageDialogOptions) (string, error)
	OpenDirectoryDialog(ctx context.Context, opts rpkg.OpenDialogOptions) (string, error)
}
