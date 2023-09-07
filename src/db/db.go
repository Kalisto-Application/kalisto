package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"kalisto/src/models"
	"os"
	"path"

	"github.com/dgraph-io/badger/v4"
)

func init() {
	gob.Register(models.Workspace{})
}

const (
	keyWorkspace  = "workspace:"
	leyGlobalVars = "globalVars"
)

type DB struct {
	db *badger.DB
}

const dbFolderName = "kalisto.db"

func New(wd string) (*DB, error) {
	dbDir := path.Join(wd, dbFolderName)
	os.MkdirAll(dbDir, os.ModePerm)
	// 65MB log file size
	db, err := badger.Open(badger.DefaultOptions(dbDir).WithValueLogFileSize(1 << 26))
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) SaveGlobalVars(vars string) error {
	return write[string](db.db, leyGlobalVars, vars)
}

func (db *DB) GlobalVars() (string, error) {
	return read[string](db.db, leyGlobalVars)
}

func (db *DB) SaveWorkspace(w models.Workspace) error {
	key := keyWorkspace + w.ID

	return write[models.Workspace](db.db, key, w)
}

func (db *DB) GetWorkspaces() ([]models.Workspace, error) {
	return readPrefix[models.Workspace](db.db, keyWorkspace)
}

func (db *DB) GetWorkspace(id string) (models.Workspace, error) {
	key := keyWorkspace + id
	return read[models.Workspace](db.db, key)
}

func (db *DB) DeleteWorkspace(id string) error {
	key := keyWorkspace + id

	return delete(db.db, key)
}

func readPrefix[T any](db *badger.DB, key string) (res []T, err error) {
	db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(key)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			if err := it.Item().Value(func(data []byte) error {
				var item T
				if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&item); err != nil {
					return fmt.Errorf("failed to decode item: %w", err)
				}

				res = append(res, item)
				return nil
			}); err != nil {
				return fmt.Errorf("failed to get item value: %w", err)
			}
		}

		return nil
	})

	return res, nil
}

func read[T any](db *badger.DB, key string) (res T, err error) {
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("failed to get %s: %w", key, err)
		}
		if err := item.Value(func(data []byte) error {
			if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func write[T any](db *badger.DB, key string, data T) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		return fmt.Errorf("failed to marshal %s: %w", key, err)
	}

	if err := db.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(key), buf.Bytes()); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to save %s to disk: %w", key, err)
	}

	return nil
}

func delete(db *badger.DB, key string) error {
	if err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	}); err != nil {
		return fmt.Errorf("failed to save %s to disk: %w", key, err)
	}

	return nil
}
