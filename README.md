# habit-service
[![CI](https://github.com/utkuufuk/habit-service/actions/workflows/ci.yml/badge.svg)](https://github.com/utkuufuk/habit-service/actions/workflows/ci.yml)

Habit service for [entrello](https://github.com/utkuufuk/entrello)

## Usage
Starting the server:
```sh
go run ./cmd/server
```

#### `GET <SERVER_URL>/entrello`
Fetch all habits that are not marked yet as done/skipped/failed.

[entrello](https://github.com/utkuufuk/entrello) will periodically call this endpoint and create a Trello card for each returned habit.

Alternatively, you can run following command to print the same result set on the console:
```sh
go run ./cmd/cli
```

#### `POST <SERVER_URL>/entrello`
Mark a habit as done/skipped/failed.

[entrello](https://github.com/utkuufuk/entrello) will call this endpoint whenever a habit card is archived on your Trello board.

#### `POST <SERVER_URL>/progress-report`
Generate and send a progress report as a Telegram message.

The `POST` request body may be empty, but you must set the `X-API-Key` header accordingly if you set the `SECRET` environment variable.

Alternatively, you can run following command to generate and send a progress report:
```sh
go run ./cmd/cli progress-report
```

## Configuration
Put your environment variables in a file called `.env`, based on `.env.example`.

| Environment Variable | Description |
|-|-|
| `TIMEZONE_LOCATION`           | Timezone, e.g. `"Europe/Istanbul"` |
| `GSHEETS_CLIENT_ID`           | Google Sheets Client ID |
| `GSHEETS_CLIENT_SECRET`       | Google Sheets Client Secret |
| `GSHEETS_ACCESS_TOKEN`        | Google Sheets Access Token |
| `GSHEETS_REFRESH_TOKEN`       | Google Sheets Refresh Token |
| `SPREADSHEET_ID`              | Google Spreadsheet ID |
| `PORT`                        | HTTP port (server mode only) |
| `SECRET`                      | API secret (server mode only, optional) |
| `PROGRESS_REPORT_SKIP_LIST`   | Comma-separated habit names to be excluded from progress reports |
| `TELEGRAM_TOKEN`              | Telegram Bot API Token |
| `TELEGRAM_CHAT_ID`            | Telegram Bot Chat ID |

## Running With Docker
A new [Docker image](https://github.com/utkuufuk?tab=packages&repo_name=habit-service) will be created upon each [release](https://github.com/utkuufuk/habit-service/releases).

1. Authenticate with the GitHub container registry (only once):
    ```sh
    echo $GITHUB_ACCESS_TOKEN | docker login ghcr.io -u GITHUB_USERNAME --password-stdin
    ```

2. Pull the latest Docker image:
    ```sh
    docker pull ghcr.io/utkuufuk/habit-service/image:latest
    ```

3. Start a container:
    ```sh
    # server
    docker run -d \
        -p <PORT>:<PORT> \
        --env-file </abs/path/to/.env> \
        --restart unless-stopped \
        --name habit-service \
        ghcr.io/utkuufuk/habit-service/image:latest

    # CLI
    docker run --rm --env-file </abs/path/to/.env> ghcr.io/utkuufuk/habit-service/image:latest ./cli
    ```
