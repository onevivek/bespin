package api

import (
	"fmt"
	"google.golang.org/api/gmail/v1"
)

func (g *Service) ListLabels(userId string) ([]*gmail.Label, error) {
	resp, err := g.s.Users.Labels.List(userId).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to list labels: %w", err)
	}

	return resp.Labels, nil
}
