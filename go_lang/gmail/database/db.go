package database

import (
	"context"

	"github.com/onevivek/bespin/go_lang/gmail/account"
)

type GetSaver interface {
	Get(ctx context.Context, id string) (*account.Account, error)
	Save(ctx context.Context, id string, a *account.Account) error
}
