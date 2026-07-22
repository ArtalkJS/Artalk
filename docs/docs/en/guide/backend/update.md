# Program Upgrade

## Update Changes

For significant changes, such as updates that require modifications to existing configurations, refer to the release notes and migration instructions:

- [v2.10.0 Release Notes and Migration Guide](../releases/v2.10.0.md)
- [Complete CHANGELOG.md](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md)

## One-Click Upgrade via Command Line

Execute `./artalk upgrade`

This operation will automatically download and upgrade the program from GitHub Release. Ensure to stop Artalk before executing this command.

::: tip
Execute `./artalk upgrade -f` with the `-f` parameter to force an update.
:::

## Docker Upgrade

Refer to: [Docker · Upgrade](./docker.md#升级)

## Manual Method

Go to [GitHub Release](https://github.com/ArtalkJS/Artalk/releases) to manually download the latest build.

Replace the old version files with the new ones.
