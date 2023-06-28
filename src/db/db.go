package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"kalisto/src/models"
	"os"
	"path"

	"github.com/peterbourgon/diskv/v3"
)

type DB struct {
	d *diskv.Diskv
}

func New(wd string) (*DB, error) {
	os.MkdirAll(path.Join(wd, "db"), os.ModePerm)

	workspaceFile := path.Join(wd, "db", "workspaces")
	if _, err := os.Stat(workspaceFile); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(workspaceFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create worksapce db: %w", err)
		}
		f.Close()
	}

	envsFile := path.Join(wd, "db", "envs")
	if _, err := os.Stat(envsFile); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(envsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create envs db: %w", err)
		}
		f.Close()
	}

	d := diskv.New(diskv.Options{
		BasePath:  path.Join(wd, "db"),
		Transform: func(s string) []string { return []string{} },
		// CacheSizeMax 10MB
		CacheSizeMax: 10 * 1024 * 1024,
	})
	return &DB{d: d}, nil
}

func (db *DB) SaveEnvs(d map[string]models.Envs) error {
	return write[map[string]models.Envs](db.d, "envs", d)
}

func (db *DB) Envs() (map[string]models.Envs, error) {
	return read[map[string]models.Envs](db.d, "envs")
}

func (db *DB) SaveWorkspaces(w []models.Workspace) error {
	return write[[]models.Workspace](db.d, "workspaces", w)
}

func (db *DB) Workspace() ([]models.Workspace, error) {
	return read[[]models.Workspace](db.d, "workspaces")
}

func read[T any](d *diskv.Diskv, key string) (data T, err error) {
	b, err := d.Read(key)
	if err != nil {
		return data, fmt.Errorf("failed to read %s from disk: %w", key, err)
	}
	if len(b) == 0 {
		return data, nil
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return data, fmt.Errorf("failed to unmarshal %s from json: %w", key, err)
	}

	return data, nil
}

func write[T any](d *diskv.Diskv, key string, data T) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal %s to json: %w", key, err)
	}
	if err := d.Write(key, b); err != nil {
		return fmt.Errorf("failed to save %s to disk: %w", key, err)
	}

	return nil
}
