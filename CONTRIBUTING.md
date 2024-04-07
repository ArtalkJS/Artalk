# Developer Contributing Guide

This guide is for developers who want to contribute to the project.

## Development Environment

To develop Artalk, including the frontend and backend, you need to install the following tools:

- [Node.js](https://nodejs.org/en/) (>= 20.12.1)
- [Go](https://golang.org/) (>= 1.22)
- [Docker](https://www.docker.com/) (>= 20.10.0) (optional)
- [Docker Compose](https://docs.docker.com/compose/) (>= 1.29.0) (optional)

Then we'll setup the development environment.

### Get the source code

Run the following command to clone the repository:

```sh
git clone https://github.com/ArtalkJS/Artalk
```

It is recommended to fork the repository first, and then clone the forked repository.

Enter the directory:

```sh
cd Artalk
```

### Build Backend

First, you need to install the dependencies for the backend written in Go. Simply run the `make install` command to install the dependencies.

Then, run the `make dev` to build and run `./bin/artalk`, and you can pass startup parameters to the program using `ARGS="version" make dev`.

This will build the backend with debugging symbols. The binary file will be placed under the `./bin` directory.

The backend program will run by default on port `23366`. You can access it through a browser at `http://localhost:23366`. It's recommended not to change this port number for testing the backend program.

### Build Frontend

First, you need to install the dependencies for the frontend. Simply run the `pnpm install` command to install the dependencies.

Then, run the `pnpm dev` to build and run the frontend, and you can pass startup parameters to the program using `ARGS="--port 5173" pnpm dev`.

The frontend program will run by default on port `5173`, and you can access it in a browser at `http://localhost:5173`. The frontend testing client will, by default, request the backend on port `23366`, so it's essential to keep the backend on this port.

The frontend program is divided into the main program and a sidebar program, with the sidebar program running on a separate port, which is `23367`.

### Development Workflow

In most cases, to set up a complete development instance, you need to navigate to the Artalk Repo directory and then execute `make dev`. This will run the backend on port `23366`. Then, execute `pnpm dev`, which will run the frontend on port `5173`. You can optionally execute `pnpm dev:sidebar` to run the sidebar frontend on port `23367`. For frontend development, you need to access `http://localhost:5173` in your browser to perform your development and testing.

When you make changes to the frontend code, you can build the complete frontend program using `pnpm build:all`. The JavaScript and CSS code can be found in `ui/artalk/dist`.

When you make changes to the backend code, running `make all` will build the complete backend program. Note that since the backend program also embeds the frontend code, the `scripts/build-frontend.sh` script will run during backend program building, which includes the embedded frontend main program and sidebar frontend program. If you are interested, you can explore the complete frontend build process in the `Makefile` code.

Additionally, there is automated CI on GitHub for building. You can find the relevant code in the `.github/workflows` directory.

### Optional: Use a One-Key Script to Run a Demo Site

When you access `http://localhost:5173` during frontend development, you will get a minimalist interface that only contains the Artalk program interface and not a real blog environment. If you want to make testing closer to a real environment, you can create a demo blog site using the following steps.

Run the following command:

```sh
./scripts/setup-example-site.sh
```

This script sets up a local example site for testing at `data` folder, with Artalk integrated into its theme.

After running this script, run:

```sh
./bin/artalk server -c ./data/artalk.yml
```

to start the artalk server.

And open <http://localhost:1313/> in your browser to view the example site.

Here is the default admin account (only created in test mode):

```yaml
name: "admin"
email: "admin@test.com"
password: "admin"
```

## Testing

The backend Go program can be tested by running `make test` for unit testing. The test results will be outputted in the terminal. You can also execute `test-coverage` to check the code test coverage.

Frontend testing is conducted using Playwright for end-to-end (E2E) testing. To start the E2E testing, run `make test-frontend-e2e`.

Both frontend and backend testing are automated and will be performed during Git pull requests and as part of the build process using CircleCI or GitHub Actions.

## Project Structure

Artalk is a monorepo project, which means all the source code is in the same repository. However, the frontend and backend are separated. The frontend part is located in `./ui` directory.

- `bin/` - The compiled binary files. This directory is ignored by git.
- `cmd/` - The source code for the command line tools.
- `conf/` - The sample configuration files.
- `docs/` - The documentation site source code.
- `i18n/` - The translation files.
- `internal/` - The internal packages.
- `data/` - The local data. This directory is ignored by git.
- `public/` - The static files. Built frontend files will be copied here.
- `scripts/` - The scripts for development.
- `server/` - The source code for the server.
- `ui/` - The source code for the frontend.

## Translation

Artalk aims to be a multilingual project. If you would like to contribute to the translation, here are some tips:

If you just write some new features or do some fixes/refactoring, use the following command to generate the translation template

```sh
go run ./internal/i18n/gen -w . -d .
```

If you're not a programmer and would like to help us improve the translation, you can edit the translation files directly in the `.i18n' directory and then submit a pull request.
