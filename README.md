# habit-service
[![CI](https://github.com/utkuufuk/habit-service/actions/workflows/ci.yml/badge.svg)](https://github.com/utkuufuk/habit-service/actions/workflows/ci.yml)

Habit service for [entrello](https://github.com/utkuufuk/entrello)

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
    # habit service
    docker run -d \
        -p <PORT>:<PORT> \
        --env-file </absolute/path/to/.env> \
        --restart unless-stopped \
        --name habit-service \
        ghcr.io/utkuufuk/habit-service/image:latest

    # score update runner
    docker run --rm \
        --env-file </absolute/path/to/.env> \
        ghcr.io/utkuufuk/habit-service/image:latest \
        ./score-update

    # progress report runner
    docker run --rm \
        --env-file </absolute/path/to/.env> \
        ghcr.io/utkuufuk/habit-service/image:latest \
        ./progress-report
    ```
