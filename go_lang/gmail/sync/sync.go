package sync

import (
	"context"
	"fmt"
	"log"

	"github.com/onevivek/bespin/go_lang/gmail/account"
	"github.com/onevivek/bespin/go_lang/gmail/api"
	"google.golang.org/api/gmail/v1"
)

type Service struct {
	sv  api.HistoryFetcher
	log Logger
}

type Logger interface {
	Printf(format string, args ...interface{})
}

func NewService(sv api.HistoryFetcher, logger Logger) *Service {
	return &Service{
		sv:  sv,
		log: logger,
	}
}

func (s *Service) Run(ctx context.Context, account account.Account, now, historicSyncDays int64) (account.Account, error) {
	s.log.Printf("started sync")
	defer s.log.Printf("stopped sync")

	if done := checkLastRun(account.LastRunAt, now); done {
		return account, nil
	}

	msgs, err := s.fetchMessages(account.ProviderCursorId, historicSyncDays)
	if err != nil {
		return account, err
	}
	s.log.Printf("fetched %d messages", len(msgs))

	account.ProviderCursorId = latestHistoryId(msgs)
	account.LastRunAt = now
	account.Status = "done"

	return account, nil
}

func (s *Service) fetchMessages(historyId uint64, historicSyncDays int64) ([]*gmail.Message, error) {
	if historyId == 0 && historicSyncDays <= 0 {
		return nil, fmt.Errorf("both historyId in account object and HISTORIC_SYNC_DAYS env are 0")
	}

	var msgs []*gmail.Message
	var err error

	if historicSyncDays > 0 {
		msgs, err = s.sv.FetchLastNDays(historicSyncDays)
	} else {
		msgs, err = s.sv.FetchStartingFromId(historyId)
	}

	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %w", err)
	}

	return msgs, nil
}

func checkLastRun(lastRunAt, now int64) bool {
	if lastRunAt > now {
		log.Println("account lastRunAt is greater than current timestamp; exiting")
		return true
	}

	return false
}

func latestHistoryId(msgs []*gmail.Message) uint64 {
	var latest uint64

	for _, msg := range msgs {
		if msg.HistoryId > latest {
			latest = msg.HistoryId
		}
	}

	return latest
}
