package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"tgBot/lib/e"
	"tgBot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' form '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownComand)

	}

}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

func (p *Processor) savePage(chatID int, pageURl string, username string) error {
	page := &storage.Page{
		URl:      pageURl,
		UserName: username,
	}

	isExist, err := p.storage.IsExists(page)
	if err != nil {
		return e.Wrap("cant do command save page", err)
	}
	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return e.Wrap("cant do command save page", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return e.Wrap("cant do command save page", err)
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return e.Wrap("cant do command: cant send random", err)
	}
	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URl); err != nil {
		return e.Wrap("cant do command: cant send random", err)
	}

	p.storage.Remove(page)

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
