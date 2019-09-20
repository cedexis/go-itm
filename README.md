[![Build Status](https://travis-ci.org/cedexis/go-itm.svg)](https://travis-ci.org/cedexis/go-itm)

# go-itm

A Go client library for accessing the Citrix ITM API.

| NOTICE: Due to strategic changes at Citrix, this project is no longer maintained. |
| --- |

## Running Unit Tests in Docker

go-itm can be updated and tested in isolation on your local machine. A Dockerfile and Make targets are provided to aid in this process.

Build a Docker image called go-itm:latest:

```bash
$ make docker-build
```

Run an interactive Bash session in a new Docker container based on the go-itm:latest image:

```bash
$ make docker-run
```

Run unit tests within the container:

```bash
[container] /go-itm $ make test 
go test ./...
ok github.com/cedexis/go-itm/itm 0.009s
```

The /go-itm directory within the container is mounted to the project root directory on the Docker host, so you can iteratively edit code on the host using your favorite editor and then re-run the unit tests inside the container.
