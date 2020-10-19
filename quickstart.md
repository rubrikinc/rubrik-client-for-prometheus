# Rubrik Prometheus Client Quick Start Guide

## Installation

### Building the Agent

The server building the agent will need GoLang installed (the package was tested with the following version: `go version go1.11 linux/amd64`).

Pull down the following dependencies:

```bash
go get github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm
go get github.com/prometheus/client_golang/prometheus
```

Clone this repository to the machine configured with GoLang, browse to the `src/golang` folder, and run the following command to build the package:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

This will build the package for the linux/amd64 architecture. For other architectures, replace the values of `GOOS` and `GOARCH` as described [here](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63).

This results in an executable named `main` in the current folder. This can be run to start exposing metrics.

### Building and running a docker image

Some users may wish to run the agent as a docker container. A Dockerfile is included in the `src/golang` folder for this purpose. To build the docker image, run the following command:

```bash
docker build -t rubrikinc/prometheus-client -f Dockerfile .
```

The resulting docker image will be in the local repository on the server.

## Using the Prometheus Agent

Ensure that the following environment variables exist, and are defined: `rubrik_cdm_node_ip`, `rubrik_cdm_username`, `rubrik_cdm_password`.

### Running from the GoLang binary

If running from the compiled GoLang binary, then we can just run `./main` to start exposing metrics on port 8080, these will then be browsable via `http://localhost:8080/metrics`.

### Running from the docker image

In the case that we are running the agent from a docker image, we can run the following command:

```bash
docker run -d -t -e rubrik_cdm_node_ip=$rubrik_cdm_node_ip \
-e rubrik_cdm_username=$rubrik_cdm_username \
-e rubrik_cdm_password=$rubrik_cdm_password \
-p 8080:8080 rubrikinc/prometheus-client
```

This will map port 8080 inside the container, to port 8080 on the docker host. Metrics will then be browsable via `http://localhost:8080/metrics`.
