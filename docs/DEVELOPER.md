# Developer Contributing Guide

This guide is for developers who want to contribute to the project.

## Development Environment

To develop Artalk, including the frontend and backend, you need to install the following tools:

- [Node.js](https://nodejs.org/en/) (>= 16.0.0)
- [Go](https://golang.org/) (>= 1.19)
- [Docker](https://www.docker.com/) (>= 20.10.0) (optional)
- [Docker Compose](https://docs.docker.com/compose/) (>= 1.29.0) (optional)

Then we'll setup the development environment.

### Get the source code

Run the following command to clone the repository:

```sh
$ git clone https://github.com/ArtalkJS/Artalk
```

It is recommended to fork the repository first, and then clone the forked repository.

Enter the directory:

```sh
$ cd Artalk
```


### Build frontend and backend

First, we need to install the dependencies for backend written in Go. Simply run the following command:

```sh
make build
```

This will build both the frontend and backend.

+ **Frontend** will be built under `./ui/packages/artalk` and copied to `./public` directory.
+ **Backend** will be built under `./bin` directory.

### Optional: Use one-key script to run a demo site

If you want to run a demo site, you can use the following command:

```sh
./scripts/setup-example-site.sh
```

This script sets up a local example site for testing, with Artalk integrated.

After running this script, run:

```sh
./bin/artalk server -c ./local/local.yml
```

to start the artalk server.

And open http://localhost:1313/ in your browser to view the example site.
