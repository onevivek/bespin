package api

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Service struct {
	s               *gmail.Service
	token           oauth2.TokenSource
	userId          string
	timeout         int
	disableBatching bool
}

// New creates new read-only api client
func New(accessToken string, timeout int, disableBatching bool) (*Service, error) {
	tokenSource, err := newTokenSource(accessToken)
	if err != nil {
		return nil, fmt.Errorf("error fetching tokenSource: %w", err)
	}

	srv, err := gmail.NewService(context.TODO(), option.WithTokenSource(tokenSource), option.WithScopes(gmail.GmailReadonlyScope))
	if err != nil {
		return nil, fmt.Errorf("failed to create new api s: %w", err)
	}

	return &Service{
		s:               srv,
		token:           tokenSource,
		userId:          "me", //TODO: maybe read this from env?
		timeout:         timeout,
		disableBatching: disableBatching,
	}, nil
}

// newTokenSource created a new oauth token source
func newTokenSource(accessToken string) (oauth2.TokenSource, error) {
	config := oauth2.Config{
		Endpoint: google.Endpoint,
	}

	token := oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	return config.TokenSource(context.TODO(), &token), nil
}

func (g *Service) GetMessage(id string) (*gmail.Message, error) {
	call := g.s.Users.Messages.Get(g.userId, id)
	res, err := call.Format("full").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve message id: %s for user: %s: %w", id, g.userId, err)
	}

	return res, nil
}

// GetLatestMessage fetches latest message with specified labels and for user. If there is no messages the returned
// message will be nil.
func (g *Service) GetLatestMessage(labelIds ...string) (*gmail.Message, error) {
	msgsId, err := g.ListLatestMessageIds(1, labelIds...)
	if err != nil {
		return nil, fmt.Errorf("issue listing message ids while getting latest message: %w", err)
	}

	if len(msgsId) == 0 {
		return nil, nil
	}

	msg, err := g.GetMessage(msgsId[0])
	if err != nil {
		return nil, fmt.Errorf("issue getting latest message: %w", err)
	}

	return msg, nil
}

func (g *Service) ListLatestMessageIds(maxResults int64, labelIds ...string) ([]string, error) {
	resp, err := g.s.Users.Messages.List(g.userId).LabelIds(labelIds...).MaxResults(maxResults).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to list messages for user %s with labels %v: %w", g.userId, labelIds, err)
	}

	var msgIds []string
	for _, msg := range resp.Messages {
		msgIds = append(msgIds, msg.Id)
	}

	return msgIds, nil
}
