package mime

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestParseMultipartBody(t *testing.T) {
	boundary := "batch_gc9oYFNRLsvqkjNwDmTKYO7n0seqGHMZ"
	body, err := ioutil.ReadFile("../examples/multipart.txt")
	if err != nil {
		t.Fatalf("error reading multipart example file: %v", err)
	}

	msgs, err := parseMultipartBody(bytes.NewReader(body), boundary)
	if err != nil {
		t.Fatalf("error parsing multipart body: %v", err)
	}

	want := 3
	got := len(msgs)
	if got != want {
		t.Fatalf("wanted messages length %d, got %d", want, got)
	}
}
