package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tgBot/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPage = errors.New("no saved pages")

type Page struct {
	URl      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URl); err != nil {
		return "", e.Wrap("cant calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("cant calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
