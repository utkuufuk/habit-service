package habit

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

const (
	AUTH_URL     = "https://accounts.google.com/o/oauth2/auth"
	REDIRECT_URL = "urn:ietf:wg:oauth:2.0:oob"
	SCOPE        = "https://www.googleapis.com/auth/spreadsheets"
	TOKEN_URL    = "https://oauth2.googleapis.com/token"
)

// initializeService initializes the gsheets service
func initializeService(ctx context.Context) (service *sheets.Service, err error) {
	config, token, err := readCreds("token.json")
	if err != nil {
		return service, fmt.Errorf("could not get credentials for google spreadsheets: %w", err)
	}

	client := config.Client(ctx, token)
	return sheets.New(client)
}

func readCreds(tokenFile string) (*oauth2.Config, *oauth2.Token, error) {
	cfg := &oauth2.Config{
		ClientID:     config.GoogleClientId,
		ClientSecret: config.GoogleClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AUTH_URL,
			TokenURL: TOKEN_URL,
		},
		RedirectURL: REDIRECT_URL,
		Scopes:      []string{SCOPE},
	}

	expiry, err := time.Parse(time.RFC3339, "2021-02-28T17:29:48.024495+03:00")
	token := &oauth2.Token{
		AccessToken:  config.GoogleAccessToken,
		TokenType:    "Bearer",
		RefreshToken: config.GoogleRefreshToken,
		Expiry:       expiry,
	}
	return cfg, token, err
}
