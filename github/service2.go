package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v69/github"
	"github.com/joakimen/gg"
)

type Service2 struct {
	Client *github.Client
}

func NewService2(authToken string) Service2 {
	timeoutSeconds := 10
	httpClient := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	return Service2{
		Client: github.NewClient(httpClient).WithAuthToken(authToken),
	}
}

func (s *Service2) GetAuthenticatedUser(ctx context.Context) (gg.GitHubUser, error) {
	user, _, err := s.Client.Users.Get(ctx, "")
	if err != nil {
		return gg.GitHubUser{}, fmt.Errorf("failed to get the authenticated user: %w", err)
	}

	mappedUser := gg.GitHubUser{
		Login: user.GetLogin(),
	}
	return mappedUser, nil
}

// Demo function
// func
