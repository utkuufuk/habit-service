## Running with Docker
The Docker image is updated whenever a new commit is pushed to the `master` branch.

*See `.github/workflows/release.yml` for continuous delivery workflow configuration.*

1. [Authenticate with GitHub Container Registry](https://docs.github.com/en/free-pro-team@latest/packages/guides/configuring-docker-for-use-with-github-packages#authenticating-to-github-packages) (only once)
    ```sh
    docker login https://docker.pkg.github.com -u utkuufuk
    ```

2. Pull the docker image
    ```sh
    docker pull docker.pkg.github.com/utkuufuk/habit-service/habit-service-image:latest
    ```

3. Create a `config.yml` in the same structure as `config.example.yml`, as well as `credentials.json` and `token.json` files for Google Sheets authentication. The following steps assumes that these files are located at `~/.config/habit-service/`.

4. Run the image
    ```sh
    docker run -d \
        -v ~/.config/habit-service/config.yml:/src/config.yml \
        -v ~/.config/habit-service/token.json:/src/token.json \
        -v ~/.config/habit-service/credentials.json:/src/credentials.json \
        -p <PORT>:<PORT> \
        --restart unless-stopped \
        --name habit_service \
        docker.pkg.github.com/utkuufuk/habit-service/habit-service-image:latest
    ```
