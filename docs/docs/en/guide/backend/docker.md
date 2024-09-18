# Docker

Artalk provides Docker images for the backend program to accelerate the deployment process and offer a good deployment experience.

The [Docker Hub](https://hub.docker.com/r/artalk/artalk-go) image versions are kept in sync with the [Releases](https://github.com/ArtalkJS/Artalk/releases) in the code repository.

## Pulling the Image

```bash
docker pull artalk/artalk-go
```

## Creating the Container

::: tip

It is recommended to use Docker Compose. Refer to the steps on the [Program Deployment](../deploy) page.

:::

## Generating the Configuration File

```bash
docker run -it -v $(pwd)/data:/data --rm artalk/artalk-go gen config data/artalk.yml
```

## Restart

After modifying the configuration file, you need to restart to apply the changes.

```bash
# Docker Compose
docker-compose restart

# Docker
docker restart artalk
```

## Stop

```bash
# Docker Compose
docker-compose stop

# Docker
docker stop artalk
```

## Upgrade

Delete the existing container, pull the latest image, and then recreate the container.

### Docker Compose

```bash
docker-compose down
docker-compose pull
docker-compose up -d
```

### Docker

```bash
docker stop artalk
docker rm artalk
docker pull artalk/artalk-go
```

::: tip
For breaking changes, please refer to the update notes in the version changelog: [CHANGELOG.md](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md)
:::

## Testing Version

The nightly image is a testing version that is updated daily and automatically built from the latest code in the repository.

```bash
docker pull artalk/artalk-go:nightly
```

## Historical Versions

Images are automatically built and released with the repository tags. You can pull images of different versions.

```bash
docker pull artalk/artalk-go:<version>
```

## Entering the Container

```bash
# Docker Compose
docker-compose exec artalk bash

# Docker
docker exec -it artalk bash
```

## Multi-Platform Compatibility

The Docker image currently only provides builds for x86 and arm64 architectures. For more platform architectures, please download the binary build and deploy it using the [Binary Deployment](../deploy.md#binary) method.
