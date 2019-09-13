# Rubrik Prometheus Client

## Summary

A Prometheus client for Rubrik CDM data point exposure, written in Go.

## Usage

Ensure that the following environment variables exist, and are defined: `rubrik_cdm_node_ip`, `rubrik_cdm_node_username`, `rubrik_cdm_node_password`.

Copy the client from `src/golang/prometheus_client.go` in this repository, and execute it using `go run .\prometheus_client.go`.

This will expose the cluster metrics on port 8080, these will then be browsable via `http://localhost:8080/metrics`.

![](https://www.saashub.com/images/app/service_logos/6/f4fb68e43ee1/large.png?1526640038)