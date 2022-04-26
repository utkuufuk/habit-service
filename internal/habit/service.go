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
	authUrl     = "https://accounts.google.com/o/oauth2/auth"
	redirectUrl = "urn:ietf:wg:oauth:2.0:oob"
	scope       = "https://www.googleapis.com/auth/spreadsheets"
	tokenUrl    = "https://oauth2.googleapis.com/token"
)

func initService(ctx context.Context, cfg config.GoogleSheetsConfig) (service *sheets.Service, err error) {
	auth := &oauth2.Config{
		ClientID:     cfg.GoogleClientId,
		ClientSecret: cfg.GoogleClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authUrl,
			TokenURL: tokenUrl,
		},
		RedirectURL: redirectUrl,
		Scopes:      []string{scope},
	}

	token := &oauth2.Token{
		AccessToken:  cfg.GoogleAccessToken,
		TokenType:    "Bearer",
		RefreshToken: cfg.GoogleRefreshToken,
		Expiry:       time.Now(),
	}
	if err != nil {
		return service, fmt.Errorf("could not get credentials for google spreadsheets: %w", err)
	}

	client := auth.Client(ctx, token)
	return sheets.New(client)
}
