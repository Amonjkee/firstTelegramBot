package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"tgBot/lib/e"
	"tgBot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return e.Wrap("cant save", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap("cant save", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap("cant save", err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap("cant save", err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	err := s.IsUserExists(userName)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)

	if err != nil {
		return nil, e.Wrap("cant pick random page", err)
	}

	if len(files) == 0 {
		return nil, e.Wrap("cant pick random page", storage.ErrNoSavedPage)
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) IsUserExists(userName string) error {
	files, _ := os.ReadDir(s.basePath)

	for _, file := range files {
		if file.Name() == userName {
			return nil
		}
	}
	return storage.ErrNoSavedPage
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("cant remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("cant remove file: %s", path)
		return e.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("cant check file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("cant check file %s exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap("cant decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("cant decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
