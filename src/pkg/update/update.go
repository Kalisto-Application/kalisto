package update

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/minio/selfupdate"
	"golang.org/x/mod/semver"
)

type Tag struct {
	Name string `json:"name"`
}

func NewUpdater(current, platform, token string) *Update {
	return &Update{
		current:  current,
		token:    token,
		platform: platform,
	}
}

type Update struct {
	current  string
	token    string
	platform string
}

func (u *Update) Run() (bool, error) {
	if u.current == "" {
		return false, nil
	}
	r, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/Kalisto-Application/kalisto-app/tags", nil)
	if err != nil {
		return false, err
	}

	r.Header.Set("Accept", "application/vnd.github+json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.token))
	r.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return false, err
	}

	if len(tags) == 0 {
		return false, fmt.Errorf("found 0 tags")
	}

	tag := tags[0].Name

	if semver.Compare(tag, u.current) <= 0 {
		return false, nil
	}

	resp, err = http.DefaultClient.Get(fmt.Sprintf("https://kzmjuiampaqvqikqsbvp.supabase.co/storage/v1/object/public/release/kalisto-bin-%s-%s.zip", u.platform, tag))
	if err != nil {
		return false, fmt.Errorf("failed to fetch new bin: %w", err)
	}
	f, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return false, fmt.Errorf("failed to create bin tmp")
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return false, fmt.Errorf("failed to copy response to a tmp zip file: %w", err)
	}

	bin, err := u.unzip(f.Name())
	if err != nil {
		return false, fmt.Errorf("failed to unzip bin: %w", err)
	}
	if err != nil {
		return false, fmt.Errorf("failed to read next untar: %w", err)
	}

	if err := selfupdate.Apply(bin, selfupdate.Options{}); err != nil {
		if rollbackErr := selfupdate.RollbackError(err); rollbackErr != nil {
			return false, rollbackErr
		}
		return false, err
	}

	return true, nil
}

func (u *Update) unzip(name string) (io.Reader, error) {
	zipReader, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	if len(zipReader.File) == 0 {
		return nil, fmt.Errorf("no files in zip")
	}
	f, err := zipReader.File[0].Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open a zip file: %w", err)
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, f); err != nil {
		return nil, fmt.Errorf("failed to copy unziped bin data: %w", err)
	}

	return buf, nil
}

func (u *Update) Restart() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	args := os.Args
	env := os.Environ()
	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err := cmd.Start()
		if err == nil {
			os.Exit(0)
		}
		return err
	}
	return syscall.Exec(self, args, env)
}
