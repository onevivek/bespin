package sync

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/onevivek/bespin/go_lang/gmail/account"
	"google.golang.org/api/gmail/v1"
)

func TestLatest(t *testing.T) {
	cases := []struct {
		name string

		input account.Account
		want  account.Account

		now              int64
		historicSyncDays int64

		fetchFromId    fetchFromId
		fetchLastNDays fetchLastNDays

		err bool
	}{
		{
			name:  "account lastRunAt is greater than now",
			now:   100,
			input: account.Account{LastRunAt: 101},
			want:  account.Account{LastRunAt: 101},
			err:   false,
		},
		{
			name: "historyId and historicSyncDays not set",
			now:  100,
			err:  true,
		},
		{
			name:             "historicSyncDays is negative",
			now:              100,
			historicSyncDays: -1,
			err:              true,
		},
		{
			name:             "fetchLastNDays returns an error",
			now:              100,
			historicSyncDays: 10,

			fetchLastNDays: func(days int64) ([]*gmail.Message, error) {
				return nil, fmt.Errorf("error")
			},

			err: true,
		},

		{
			name: "fetchFromId returns an error",
			now:  100,

			fetchFromId: func(startHistoryId uint64) ([]*gmail.Message, error) {
				return nil, fmt.Errorf("error")
			},

			err: true,
		},

		{
			name:  "history returned single response",
			now:   100,
			input: account.Account{ProviderCursorId: 1233},
			want:  account.Account{LastRunAt: 100, ProviderCursorId: 1234, Status: "done"},

			fetchFromId: func(startHistoryId uint64) ([]*gmail.Message, error) {
				return []*gmail.Message{{HistoryId: 1234}}, nil
			},

			err: false,
		},

		{
			name:  "history returned multiple responses",
			now:   100,
			input: account.Account{ProviderCursorId: 1233},
			want:  account.Account{LastRunAt: 100, ProviderCursorId: 1235, Status: "done"},

			fetchFromId: func(startHistoryId uint64) ([]*gmail.Message, error) {
				return []*gmail.Message{{HistoryId: 1234}, {HistoryId: 1235}}, nil
			},

			err: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sv := NewMockService(c.fetchFromId, c.fetchLastNDays)
			syncService := NewService(sv, &MockLogger{})
			got, err := syncService.Run(context.TODO(), c.input, c.now, c.historicSyncDays)
			if c.err && err == nil {
				t.Fatal("expected an error, got nil")
			}

			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("want %v, got %v\n", c.want, got)
			}
		})
	}
}

type MockLogger struct{}

func (m *MockLogger) Printf(format string, args ...interface{}) {}

//--------------------------------

type fetchFromId func(startHistoryId uint64) ([]*gmail.Message, error)
type fetchLastNDays func(days int64) ([]*gmail.Message, error)

type MockService struct {
	fetchFromId    fetchFromId
	fetchLastNDays fetchLastNDays
}

func NewMockService(fetchFromId fetchFromId, fetchLastNDays fetchLastNDays) *MockService {
	return &MockService{
		fetchFromId:    fetchFromId,
		fetchLastNDays: fetchLastNDays,
	}
}

func (m *MockService) FetchStartingFromId(startHistoryId uint64) ([]*gmail.Message, error) {
	return m.fetchFromId(startHistoryId)
}

func (m *MockService) FetchLastNDays(days int64) ([]*gmail.Message, error) {
	return m.fetchLastNDays(days)
}
