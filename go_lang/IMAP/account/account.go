package account

import (
	"fmt"
	"strconv"
)

type Account struct {
	LastRunAt        int64  `json:"last_run_at"`
	ProviderCursorId uint64 `json:"provider_cursor_id"`
	Status           string `json:"status"`
}

func ParseMap(h map[string]string) (Account, error) {
	var lastRunAt int64
	var providerCursorId uint64
	var err error

	if lastRunAtStr, ok := h["last_run_at"]; ok {
		if lastRunAt, err = strconv.ParseInt(lastRunAtStr, 10, 64); err != nil {
			return Account{}, fmt.Errorf("error parsing 'lastRunAt' %s", lastRunAtStr)
		}
	}

	if providerCursorIdStr, ok := h["provider_cursor_id"]; ok {
		if providerCursorId, err = strconv.ParseUint(providerCursorIdStr, 10, 64); err != nil {
			return Account{}, fmt.Errorf("error parsing 'providerCursorId' %s", providerCursorIdStr)
		}
	}

	return Account{
		LastRunAt:        lastRunAt,
		ProviderCursorId: providerCursorId,
	}, nil
}

func ParseAccount(a Account) map[string]interface{} {
	result := make(map[string]interface{})
	result["last_run_at"] = strconv.FormatInt(a.LastRunAt, 10)
	result["provider_cursor_id"] = strconv.FormatUint(a.ProviderCursorId, 10)
	result["status"] = a.Status

	return result
}
