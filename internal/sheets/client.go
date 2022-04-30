package sheets

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

type Client struct {
	spreadsheetId string
	service       *sheets.SpreadsheetsValuesService
}

func GetClient(ctx context.Context, cfg config.GoogleSheetsConfig) (client Client, err error) {
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

	service, err := sheets.New(auth.Client(ctx, token))
	if err != nil {
		return client, fmt.Errorf("could not initialize gsheets service: %w", err)
	}

	return Client{cfg.SpreadsheetId, service.Spreadsheets.Values}, nil
}
