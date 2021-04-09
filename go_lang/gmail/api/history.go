package api

import (
	"fmt"
	"time"

	"google.golang.org/api/gmail/v1"
)

const MaxResults = 100

type HistoryFetcher interface {
	FetchStartingFromId(startHistoryId uint64) ([]*gmail.Message, error)
	FetchLastNDays(days int64) ([]*gmail.Message, error)
}

func (g *Service) FetchStartingFromId(startHistoryId uint64) ([]*gmail.Message, error) {
	var msgIds []string

	resp, nextPageToken, err := g.fetchPageOfHistoryStartingFromId(startHistoryId, "")
	if err != nil {
		return nil, err
	}
	msgIds = append(msgIds, resp...)

	for nextPageToken != "" {
		resp, nextPageToken, err = g.fetchPageOfHistoryStartingFromId(startHistoryId, nextPageToken)
		if err != nil {
			return nil, err
		}

		msgIds = append(msgIds, resp...)
	}

	msgs, err := g.GetMessages(msgIds)
	if err != nil {
		return nil, fmt.Errorf("error getting messages: %w", err)
	}

	return msgs, nil
}

func (g *Service) fetchPageOfHistoryStartingFromId(startHistoryId uint64, nextPageToken string) ([]string, string, error) {
	req := g.s.Users.History.List(g.userId).StartHistoryId(startHistoryId).MaxResults(MaxResults)
	if nextPageToken != "" {
		req = req.PageToken(nextPageToken)
	}

	resp, err := req.Do()
	if err != nil {
		return nil, "", fmt.Errorf("unable to list page of history: %w", err)
	}

	return extractMsgIdsFromHistory(resp.History), resp.NextPageToken, nil
}

func extractMsgIdsFromHistory(history []*gmail.History) []string {
	var result []string

	for _, h := range history {
		for _, m := range h.Messages {
			result = append(result, m.Id)
		}
	}

	return result
}

func (g *Service) FetchLastNDays(days int64) ([]*gmail.Message, error) {
	now := time.Now().Unix() * 1000

	msgIds, err := g.fetchLastNDaysOfMessageIds(days, now)
	if err != nil {
		return nil, err
	}

	msgs, err := g.GetMessages(msgIds)
	if err != nil {
		return nil, fmt.Errorf("error getting messages: %w", err)
	}

	return msgs, nil
}

func (g *Service) fetchLastNDaysOfMessageIds(days int64, now int64) ([]string, error) {
	var messageIds []string

	msgIds, page, lastInternalDate, err := g.fetchPageOfMessageIds("")
	if err != nil {
		return nil, err
	}

	messageIds = append(messageIds, msgIds...)

	if page == "" {
		return messageIds, nil
	}

	daysInMiliseconds := days * 24 * 3600 * 1000
	for now-lastInternalDate < daysInMiliseconds {
		msgIds, page, lastInternalDate, err = g.fetchPageOfMessageIds(page)
		if err != nil {
			return nil, err
		}

		messageIds = append(messageIds, msgIds...)

		if page == "" {
			break
		}
	}

	return messageIds, nil
}

// fetchPageOfMessageIds takes a page token and returns message ids, next page token, internal date of the last message
// and an error
func (g *Service) fetchPageOfMessageIds(pageToken string) ([]string, string, int64, error) {
	msgs, nextPage, err := g.fetchPageOfMessages(pageToken)
	if err != nil {
		return nil, "", 0, fmt.Errorf("error fetching page of messages: %w", err)
	}

	if len(msgs) == 0 {
		return []string{}, "", 0, nil
	}

	var msgIds []string
	for _, msg := range msgs {
		msgIds = append(msgIds, msg.Id)
	}

	if nextPage == "" {
		return msgIds, "", 0, nil
	}

	lastId := msgs[len(msgs)-1].Id // TODO: seems that gmail.List returns messages in chronological order, but we should test that
	last, err := g.GetMessage(lastId)
	if err != nil {
		return nil, "", 0, fmt.Errorf("error fetching message %s: %w", lastId, err)
	}

	return msgIds, nextPage, last.InternalDate, nil
}

func (g *Service) fetchPageOfMessages(pageToken string) ([]*gmail.Message, string, error) {
	req := g.s.Users.Messages.List(g.userId).MaxResults(MaxResults)
	if pageToken != "" {
		req = req.PageToken(pageToken)
	}

	resp, err := req.Do()
	if err != nil {
		return nil, "", fmt.Errorf("unable to list messages for user %s: %w", g.userId, err)
	}

	if len(resp.Messages) == 0 {
		return nil, "", nil
	}

	return resp.Messages, resp.NextPageToken, nil
}
