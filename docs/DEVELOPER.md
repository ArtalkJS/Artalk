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


### Setup a backend for development

First, we need to install the dependencies for backend written in Go. Simply run the following command:

```sh

