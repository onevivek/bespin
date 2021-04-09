package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/onevivek/bespin/go_lang/gmail/kit"
	"github.com/onevivek/bespin/go_lang/gmail/mime"
	"google.golang.org/api/gmail/v1"
)

// https://developers.google.com/gmail/api/guides/batch
const MaxCallsPerSingleBatch = 100

func (g *Service) GetMessages(messageIds []string) ([]*gmail.Message, error) {
	if len(messageIds) == 0 {
		return []*gmail.Message{}, nil
	}

	var msgs []*gmail.Message
	var err error
	if g.disableBatching {
		msgs, err = g.getMessages(messageIds)
	} else {
		msgs, err = g.batchGetMessages(messageIds)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting messages: %w", err)
	}

	return msgs, nil
}

func (g *Service) getMessages(messageIds []string) ([]*gmail.Message, error) {
	result := make([]*gmail.Message, 0, len(messageIds))
	for _, msgId := range messageIds {
		msg, err := g.GetMessage(msgId)
		if err != nil {
			return nil, fmt.Errorf("error fetching message %s: %v", msgId, err)
		}

		result = append(result, msg)
	}

	return result, nil
}

func (g *Service) batchGetMessages(messageIds []string) ([]*gmail.Message, error) {
	client := &http.Client{Timeout: time.Duration(g.timeout) * time.Second}
	defer client.CloseIdleConnections()

	var result []*gmail.Message
	for _, msgSplit := range kit.Split(messageIds, MaxCallsPerSingleBatch) {
		singleBatchResult, err := g.singleBatchGetMessages(client, msgSplit)
		if err != nil {
			return nil, fmt.Errorf("error fetching messages batch %v", err) //TODO: add retry with backoff
		}

		result = append(result, singleBatchResult...)
	}

	return result, nil
}

func (g *Service) singleBatchGetMessages(client *http.Client, msgIds []string) ([]*gmail.Message, error) {
	req, err := g.createBatchRequest(msgIds)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with response code: %d; %s", resp.StatusCode, body)
	}

	return mime.ProcessBatchResponse(resp)
}

func (g *Service) createBatchRequest(msgIds []string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, msgId := range msgIds {
		m := textproto.MIMEHeader{
			"Content-Type": []string{"application/http"},
			"Content-Id":   []string{fmt.Sprintf("msgId:%s", msgId)},
		}

		pw, err := writer.CreatePart(m)
		if err != nil {
			return nil, fmt.Errorf("couldn't create new multipart section: %w", err)
		}

		uri := fmt.Sprintf("api.googleapis.com/gmail/v1/users/%s/messages/%s", g.userId, msgId)
		if _, err := pw.Write([]byte(fmt.Sprintf("GET %s HTTP/1.1\r\n", uri))); err != nil {
			return nil, fmt.Errorf("couldn't write multipart section: %w", err)
		}

		if _, err := pw.Write([]byte("\r\n")); err != nil {
			return nil, fmt.Errorf("couldn't write empty line in multipart section: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("issue closing multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://www.googleapis.com/batch/gmail/v1", bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("issue creating new http request: %w", err)
	}

	token, err := g.token.Token()
	if err != nil {
		return nil, fmt.Errorf("issue fetching new bearer access token: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary()))
	req.ContentLength = int64(body.Len())

	return req, nil
}
