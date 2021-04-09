package main

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nylas/go-kit/logwrapper"
	"github.com/onevivek/bespin/go_lang/gmail/api"
	"github.com/onevivek/bespin/go_lang/gmail/database"
	"github.com/onevivek/bespin/go_lang/gmail/sync"
)

type Specification struct {
	AccessToken string `required:"true" split_words:"true"`

	RedisAddress string `required:"true" split_words:"true" default:"redis:6379"`

	AccountId        string `required:"true" split_words:"true"`
	BatchTimeout     int    `split_words:"true" default:"10"`
	HistoricSyncDays int64  `split_words:"true"`
	DisableBatching  bool   `split_words:"true" default:"false"`
}

func main() {
	logger := logwrapper.NewLogger()

	var spec Specification
	err := envconfig.Process("gmail", &spec)
	if err != nil {
		logger.MissingEnvVar(err)
		return
	}

	db, err := database.NewRedis(spec.RedisAddress)
	if err != nil {
		logger.ConnectionOpenError(spec.RedisAddress, err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.ConnectionCloseError(spec.RedisAddress, err)
		}
	}()

	sv, err := api.New(spec.AccessToken, spec.BatchTimeout, spec.DisableBatching)
	if err != nil {
		logger.Errorf("cannot create api client: %v\n", err)
		return
	}

	account, err := db.Get(context.TODO(), spec.AccountId)
	if err != nil {
		logger.Errorf("cannot get account id %s from db: %v\n", spec.AccountId, err)
		return
	}

	now := time.Now().Unix()
	syncedAccount, err := sync.NewService(sv, logger).Run(context.TODO(), account, now, spec.HistoricSyncDays)
	if err != nil {
		logger.SyncError(err)
		return
	}

	if err := db.Save(context.TODO(), spec.AccountId, syncedAccount); err != nil {
		logger.Errorf("error saving account to db: %v", err)
		return
	}
}
