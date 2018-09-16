# geo-cli
Geo cli tool written in Golang using MaxMind geo data

[![GoDoc](https://godoc.org/github.com/major1201/geo-cli?status.svg)](https://godoc.org/github.com/major1201/geo-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/major1201/geo-cli)](https://goreportcard.com/report/github.com/major1201/geo-cli)

## Installation

```bash
$ go install github.com/major1201/geo-cli/cmd/geo
```

## Usage

Download a MaxMind geo data file in <https://dev.maxmind.com/geoip/geoip2/geolite2/>

Set the environment var

```bash
export GEO_MMDBFILE=/opt/GeoLite2-City.mmdb
```

Query one address in detail

```bash
geo --detail 81.2.69.142
```

Query multiple addresses in one-line format

```bash
geo www.google.com 81.2.69.142
```

Read from pipe

```bash
traceroute www.google.com | geo
```

Specify language

```bash
geo --language zh-CN 81.2.69.142
```

## Contributing

Just fork the repositry and open a pull request with your changes.

## Licence

MIT
