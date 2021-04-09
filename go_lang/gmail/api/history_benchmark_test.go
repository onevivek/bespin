package api

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func BenchmarkService_FetchLastNDays(b *testing.B) {
	var accessToken string
	var ok bool

	if accessToken, ok = os.LookupEnv("GMAIL_ACCESS_TOKEN"); !ok {
		b.Fatal("GMAIL_ACCESS_TOKEN env not set")
	}

	svc, err := New(accessToken, 0, false)
	if err != nil {
		b.Fatalf("error creating gmail service: %f", err)
	}

	days := []int64{5, 10, 20, 50, 100}

	for _, d := range days {
		startNano := time.Now().UnixNano()

		m, err := svc.FetchLastNDays(d)
		if err != nil {
			b.Errorf("error fetching last %d days: %v", d, err)
		}

		durationNano := time.Now().UnixNano() - startNano
		fmt.Printf("fetching %d messages from the last %d days took: %d nanoseconds = %d seconds\n", len(m), d, durationNano, durationNano/1_000_000_000)
	}
}
