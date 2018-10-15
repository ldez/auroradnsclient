Aurora DNS API client
=====================

[![Build Status](https://img.shields.io/travis/edeckers/auroradnsclient/master.svg?maxAge=2592000&style=flat-square)](https://travis-ci.org/edeckers/auroradnsclient)
[![License](https://img.shields.io/github/license/edeckers/auroradnsclient.svg?maxAge=2592000&style=flat-square)](https://www.mozilla.org/en-US/MPL/2.0)

An wrapper library for the Aurora DNS API, written in Go.

## Features

* List zones and records
* Add and remove records

## Requirements

* Go >= 1.6 (Go >= 1.11 for dev)

## Validate

```bash
make
```

## Basic usage

```go
tr, _ := aurora.NewTokenTransport("userID", "key")
client, _ := aurora.NewClient(tr.Client())

zones, _, _ := client.GetZones()

fmt.Println(zones)
```

## License

`auroradnsclient` is licensed under MPL-2.0 - see the LICENSE file for details
