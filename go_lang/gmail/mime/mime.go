package mime

import (
	"bufio"
	"encoding/json"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

// ProcessBatchResponse parses Gmail multipart batch response and returns slice of messages
func ProcessBatchResponse(resp *http.Response) ([]*gmail.Message, error) {
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("issue parsing response content-type: %w", err)
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, fmt.Errorf("invalid Content-Type returned %s", mediaType)
	}

	var boundary string
	var ok bool
	if boundary, ok = params["boundary"]; !ok {
		return nil, fmt.Errorf("mime boundary not set")
	}

	msgs, err := parseMultipartBody(resp.Body, boundary)
	if err != nil {
		return nil, fmt.Errorf("issue parsing multipart body: %w", err)
	}

	return msgs, nil
}

func parseMultipartBody(bodyReader io.Reader, boundary string) ([]*gmail.Message, error) {
	var result []*gmail.Message
	mr := multipart.NewReader(bodyReader, boundary)

	for {
		pr, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		m, err := parseMimePart(pr)
		if err != nil {
			return nil, err
		}

		result = append(result, m)
	}

	return result, nil
}

func parseMimePart(pr *multipart.Part) (*gmail.Message, error) {
	defer func() {
		if err := pr.Close(); err != nil {
			log.Println("error closing multipart part reader")
		}
	}()

	if pr.Header.Get("Content-Type") != "application/http" {
		return nil, fmt.Errorf("invalid Content-Type: %s while parsing Mime", pr.Header.Get("Content-Type"))
	}

	var resp *http.Response
	var err error
	if resp, err = http.ReadResponse(bufio.NewReader(pr), nil); err != nil {
		return nil, fmt.Errorf("invalid Content Type: %s", pr.Header.Get("Content-Type"))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("error closing multipart part reader")
		}
	}()

	m := gmail.Message{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode json body: %w", err)
	}

	return &m, nil
}
