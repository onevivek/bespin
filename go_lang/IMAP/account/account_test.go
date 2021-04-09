package account

import (
	"reflect"
	"testing"
)

func TestParseMap(t *testing.T) {
	cases := []struct {
		name  string
		input map[string]string
		want  Account
		err   bool
	}{
		{
			name:  "invalid lastRunAt",
			input: map[string]string{"last_run_at": "aaa"},
			err:   true,
		},
		{
			name:  "invalid provider_cursor_id",
			input: map[string]string{"provider_cursor_id": "bbb"},
			err:   true,
		},
		{
			name:  "valid account",
			input: map[string]string{"last_run_at": "123", "provider_cursor_id": "456"},
			want: Account{
				LastRunAt:        123,
				ProviderCursorId: 456,
			},
			err: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseMap(c.input)
			if c.err && err == nil {
				t.Fatalf("expected an error but got nil")
			}

			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("want %v, got %v", c.want, got)
			}
		})
	}
}
