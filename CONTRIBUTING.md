# Developer Contributing Guide

This guide helps developers understand how to contribute code, documentation, tests, and feedback to the Artalk project. We welcome all forms of contributions, whether it's large feature development, small fixes, documentation improvements, or testing. We believe that every contribution is an important part of moving the project forward.

This guide is available in multiple languages: [English](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md) • [简体中文](https://artalk.js.org/develop/contributing.html)

## Setup

Clone the repository:

```sh
git clone https://github.com/ArtalkJS/Artalk.git
```

It's recommended to fork the repository first, and then clone the forked repository.

Navigate to the directory:

```sh
cd Artalk
```

## Development Environment

To develop Artalk, both frontend and backend, install the following tools:

- [Node.js](https://nodejs.org/en/) (>= 20.17.0)
- [PNPM](https://pnpm.io/) (>= 9.10.0)
- [Go](https://golang.org/) (>= 1.22)
- [Docker](https://www.docker.com/) (>= 20.10.0) (optional)
- [Docker Compose](https://docs.docker.com/compose/) (>= 1.29.0) (optional)

## Development Workflow Overview

The development workflow consists of the following steps:

1. **Set Up Complete Development Instance**:

   - Navigate to the Artalk Repo directory.
   - Execute `make dev` to run the backend on port `23366`.
   - Execute `pnpm dev` to run the frontend on port `5173`.
   - Optionally, execute `pnpm dev:sidebar` to run the sidebar frontend on port `23367`.

2. **Access Frontend for Development**:

   - Open your browser and go to `http://localhost:5173` to perform your development and testing.

3. **Testing**:

   - Run the backend unit tests using `make test`.
   - Run the ui unit tests using `pnpm test`.
   - Run the ui E2E tests using `make test-frontend-e2e`.

4. **Building Frontend Code**:

   - When you make changes to the frontend code, build the complete frontend program using `pnpm build:all`.
   - The JavaScript and CSS code will be found in `ui/artalk/dist`.
   - The frontend code is automatically deployed to NPM by GitHub Actions.

5. **Building Backend Code**:

   - When you make changes to the backend code, note that the backend program embeds the frontend code.
   - Before building the backend, run `make build-frontend` (which runs `scripts/build-frontend.sh`). This will place the embedded frontend main program and sidebar frontend program in `/public`.
   - Then, build the backend program using `make build`.

More exploration:

- **Developing Artalk Plugins or Themes**:
  Refer to the [Artalk Plugin Development Guide](https://artalk.js.org/develop/plugin.html).

- **Makefile**:
  You can explore the process and commands in the `Makefile`.

- **Automated CI on GitHub**:
  The automated CI pipeline on GitHub, located in the `.github/workflows` directory, builds the frontend, backend, and Docker images. It publishes Docker images to Docker Hub, handles npm package publishing, and creates GitHub Releases. And also testing the frontend and backend.

The following sections provide more detailed information on the development workflow.

## Project Structure

Artalk is a monorepo project, meaning all the source code resides in a single repository. Despite this, the frontend and backend are kept separate. The frontend part is located in the `./ui` directory.

### Directory Overview

- `bin/`: The compiled binary files. (This directory is ignored by git).
- `cmd/`: The source code for the command line tools.
- `conf/`: The sample configuration files.
- `data/`: The local data. (This directory is ignored by git).
- `docs/`: The documentation site source code.
- `i18n/`: The translation files.
- `local/`: The local files. (This directory is ignored by git).
- `internal/`: The Go source code for the internal packages.
- `public/`: The static files. Built frontend files will be copied here.
- `scripts/`: The scripts for development and building.
- `server/`: The Go source code for the backend.
- `ui/`: The UI source code for the frontend.
- `.github/`: The GitHub Actions workflows.
- `.vscode/`: The VSCode settings.

## Backend Development

### Run Backend

1. **Create Configuration File**:

   - Copy `./conf/artalk.example.yml` to the root directory.
   - Rename it to `artalk.yml`.
   - Modify the file as needed.

2. **Start Backend Program**:

   - Run `make dev` to start the backend on port `23366`.
   - Access it via `http://localhost:23366`.
   - It's recommended to keep the default port for testing.
   - (Alternative) Use `ARGS="version" make dev` to pass startup parameters (default is `ARGS="server"`).

The `make dev` builds the backend with debugging symbols, which is convenient for debugging with GDB.

If you are using the VSCode, you can use the `F5` key (or RUN button) to start debugging the backend program directly.

### Build Backend

1. **Fetch Latest Git Tag**:
   Execute `git fetch --tags` to get the latest git tag information. This tag will be used as the version number of the backend app. Alternatively, specify it using environment variables: `VERSION="v1.0.0"` and `COMMIT_HASH="66128e"`.

2. **Build Frontend**:
   Run `make build-fronted` to build the frontend, as the backend app embeds the frontend code.

3. **Build Backend**:
   Execute `make build` to build the backend program.

The built binary file will be located in the `./bin` directory.

## Frontend Development (UI)

### Run Frontend

1. **Install Dependencies**:  
   Run `pnpm install` to install all necessary dependencies.

2. **Start Development Server**:  
   Execute `pnpm dev --port 5173` to start the frontend development server.

3. **Access the Frontend**:  
   The frontend application will run on port `5173` by default. Access it in your browser at `http://localhost:5173`.

4. **Sidebar Development**:  
   To run the sidebar part of the frontend, use `pnpm dev:sidebar`.

The ports `5173` and `23367` are used for frontend development and testing. These ports are included in the `ATK_TRUSTED_DOMAINS` configuration of the backend program.

### Build Frontend

To build the frontend, run the following command:

```sh
make build-frontend
```

The compiled JavaScript and CSS files will be located in the `/public` directory, embedded into the backend program.

For more details, refer to the `scripts/build-frontend.sh` script.

### Publish Frontend

The core artalk client code is automatically deployed to NPM by GitHub Actions. The deployment process is defined in the `.github/workflows/build-ui.yml` file.

If you are writing a artalk plugin in monorepo, it would not be automatically deployed. You should publish it manually. To publish the plugin, first ensure your version number is updated in the `package.json` file. Then, run the following command:

```sh
pnpm publish --access public
```

#### Check Published Versions

The `pnpm check:publish` script is designed to verify that all packages in the project have the latest versions published on npm. It skips private packages automatically and allows filtering to check specific packages using the `-F` option.

Check all packages:

```sh
pnpm check:publish
```

This will initiate the process of checking all public packages in the repository to ensure they are up to date on npm.

### Develop and Debug `@artalk/plugin-kit` and `eslint-plugin-artalk`

The monorepo will default to installing the latest version of `@artalk/plugin-kit` and `eslint-plugin-artalk` from NPM published packages. To develop and debug these packages, run the following commands from the root of the repository:

```sh
pnpm link --global --dir ui/eslint-plugin-artalk
pnpm link --global --dir ui/plugin-kit

pnpm link eslint-plugin-artalk
pnpm link @artalk/plugin-kit
```

After executing these commands, the `eslint-plugin-artalk` and `@artalk/plugin-kit` packages will be linked to the local development environment within the monorepo workspace. You can modify the source code located in the `ui/eslint-plugin-artalk` and `ui/plugin-kit` directories, and the changes will be automatically reflected in the monorepo workspace. If you wish to unlink the packages, run the following commands:

```sh
pnpm unlink eslint-plugin-artalk
pnpm unlink @artalk/plugin-kit

pnpm uninstall --global eslint-plugin-artalk
pnpm uninstall --global @artalk/plugin-kit
```

For more details, refer to the [pnpm link documentation](https://pnpm.io/cli/link).

## Docker Development

### Build Docker Image

To build the Docker image, run the following command:

```sh
docker build -t artalk:TAG .
```

Replace `TAG` with the desired tag name (e.g., `latest`).

### Skip Frontend Build

If you have already built the frontend outside the Docker container, you can skip the frontend build process inside the container to speed up the build. Use the following command:

```sh
docker build --build-arg SKIP_UI_BUILD=true -t artalk:latest .
```

For more details, refer to the `Dockerfile`.

## Testing

### Testing the Backend

To test the backend Go program, you can use the following commands:

- **Unit Testing**:
  Run `make test` to execute unit tests. The results will be outputted in the terminal.

- **Test Coverage**:
  Run `make test-coverage` to check the code test coverage.

- **HTML Coverage Report**:
  Run `make test-coverage-html` to generate and browse the test coverage report in HTML format.

- **Run Specific Tests**:
  To run specific tests, use the command with an environment variable:
  ```sh
  TEST_PATHS="./server/..." make test
  ```

### Testing the Frontend

Frontend testing includes unit tests and end-to-end (E2E) tests.

- **Unit Testing**:
  The unit tests for the frontend are conducted using Vitest. Run `pnpm test` to start the unit testing.

- **End-to-End Testing**:
  E2E testing is conducted using Playwright. To start the E2E testing, run:
  ```sh
  make test-frontend-e2e
  ```

### Continuous Integration for Testing

Both frontend and backend testing are automated and will be performed during Git pull requests and as part of the build process using CircleCI or GitHub Actions.

## Configurations

The template configuration files are located in the `/conf` directory, with templates annotated in multiple languages named in the format `artalk.example.[lang].yml` (e.g., `artalk.example.zh-CN.yml`). Artalk will parse these template configuration files to generate the settings interface, environment variable names, and documentation. To improve performance, some data will be cached during the program's runtime to eliminate parsing time.

In the `internal/config/config.go` file, there is a definition of the `Config` struct, which is used to parse the yml files into Go structs. This struct contains more precise type definitions. If you need to add, modify, or delete a configuration item, you must make changes to both the configuration file templates in `/conf` and the `Config` struct.

When you modify the configuration files, please run the following command to update the configuration data cache:

```sh
make update-conf
```

To update the environment variable documentation, you can run the command:

```sh
make update-conf-docs
```

## Documentations

The documentation consists of three parts, each associated with a specific package: **Guide documentation**, **Landing page**, and **Swagger API documentation**. Below is an organized summary of the relevant commands and their outputs for each package.

### Guide Documentation (`docs`)

- **Directory**: `docs/docs/guide`
- **Dev Command**: `pnpm -F docs dev:docs`
- **Build Command**: `pnpm -F docs build:docs`
- **Output Directory**: `docs/docs/.vitepress/dist`

### Landing Page (`docs-landing`)

- **Directory**: `docs/landing`
- **Dev Command**: `pnpm -F docs-landing dev:landing`
- **Build Command**: `pnpm -F docs-landing build:landing`
- **Output Directory**: `docs/landing/dist`

### Swagger API Documentation (`docs-swagger`)

- **Directory**: `docs/swagger`
- **Build Command**: `make update-swagger`
- **Output Directory**: `docs/swagger`

Swagger definitions are located in the backend code at `/server`. After modifying the Swagger definitions, running the build command will update the Swagger docs and generate HTTP client code at `/ui/artalk/src/api/v2.ts`.

### Combined Build

- **Command to Build All Packages**: `pnpm build:docs`
- **Combined Output Directory**: `docs/docs/.vitepress/dist`

This command builds and combines all three packages into a single directory, ready for publishing.

### Environment Variable Documentation

- **File**: `docs/docs/guide/env.md`
- **Command to Update**: `make update-conf-docs`

This command reads the configuration template file located in `/conf` to update the environment variable documentation. Run this command after modifying the configuration template file.

## Translation (i18n)

If you write new features or make fixes/refactoring, use the following command to incrementally generate the translation files by parsing the source code for `i18n.T` function calls:

```sh
make update-i18n
```

If you're not a programmer and would like to help improve the translation, you can edit the translation files directly in the `/i18n` directory and then submit a pull request.

## The End

Thank you for your patience in reading this and for your interest and support in Artalk. We understand that the success of an open-source project relies on the contributions of every developer. Whether through code, documentation, testing, or feedback, your participation is the driving force behind the project's continuous progress. We sincerely invite you to join us in improving and enhancing Artalk.

By collaborating, we can create a more robust and versatile platform that benefits the entire community. Your unique insights and expertise are invaluable, and together, we can achieve remarkable advancements. Whether you are a seasoned developer or just starting, there is always a place for your contributions in our project.

If you have any questions or need assistance, please feel free to ask in the issue section on our GitHub page: https://github.com/ArtalkJS/Artalk/issues. Thank you once again for your support, and we look forward to working with you.
