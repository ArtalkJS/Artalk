# Program Upgrade

## Update Changes

For significant changes, such as updates that require modifications to existing configurations, refer to the release notes and migration instructions:

- [v2.10.0 Release Notes and Migration Guide](../releases/v2.10.0.md)
- [Complete CHANGELOG.md](https://github.com/ArtalkJS/Artalk/blob/master/CHANGELOG.md)

## Security Branch Compatibility

These notes apply to the unreleased `security-audit-boundaries` branch.

- Automatic migration adds `guest_token_hash` to comments (default empty string) and `attempts` to email verification challenges (default 0). The earlier development column `guest_key_hash` is renamed without losing ownership, and comment caches are repopulated. Repeated migration preserves existing proofs and attempt budgets. Neither is a public comment or user API field.
- Guest credentials bind only newly created comments submitted with that credential. Names and email addresses cannot claim historical comments. Old clients can still submit comments, create pages automatically, and count page views, but private pending comments and notifications require a valid login or a guest credential held by a new client. Clearing browser storage or changing browsers does not recover a guest credential automatically.
- Artrans import/export retains its existing content format and does not transfer guest ownership proofs. To preserve those proofs during migration, use a full database backup and retain the backend address and credential storage in the browser.
- The page list parameter `has_comments` remains optional; omitting it preserves the complete list. The authentication response field `callback_origin` is also optional. With an old backend that omits it, the new authentication plugin falls back to the login entry origin. Deployments with a separate callback origin should upgrade both the backend and authentication plugin.

Older programs can ignore the added columns, but reverting to an unpatched version restores the original permission flaws. Schema compatibility does not make a security rollback safe.

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
