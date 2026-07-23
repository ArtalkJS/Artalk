# Artalk RPM packaging

This branch contains the Fedora packaging files. Go dependencies are included in the source archive. The frontend comes from the same release's `artalk_ui.tar.gz` and is checked against the SHA-256 digest from GitHub's release API. Node.js and pnpm are not required.

## COPR setup

For the `artalk` package in `@artalk/artalk`, select **Custom script** as the source type. Copy `.copr/bootstrap.sh` into the script field. Set the result directory to `results` and build dependencies to `git golang jq curl tar gzip`. Use `fedora-latest-x86_64` as the source build chroot.

The release workflow sends `{"version":"v2.10.1"}` to the Custom webhook. A version is required; the script never selects the latest tag. Automatic builds on Git pushes can stay disabled.

The build sets the RPM version from this tag. It keeps the spec's packaging revision for the same version and uses revision `1` for other versions.

## Local build

Install Go, Git, jq, curl, tar, gzip, make, rpm-build, and systemd-rpm-macros. From this directory, run:

```sh
make -f .copr/Makefile srpm outdir="$PWD/results" version=v2.10.1
rpmbuild --rebuild results/artalk-*.src.rpm
```

Use a clean Fedora build environment. The RPM build runs without network access. Test installation, service startup, the version API, and upgrades with existing configuration and comment data before publishing packaging changes.

For v2.10.1 only, the build excludes the `server/handler` test package because its background tasks can access a test database after cleanup (upstream PR #1205). Other Go tests still run. Later versions run the full suite.
