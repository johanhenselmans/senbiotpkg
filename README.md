# Senbiot

[![Build Status](https://travis-ci.org/johanhenselmans/senbiot.svg?branch=master)](https://travis-ci.org/johanhenselmans/senbiot)
[![Go Report Card](https://goreportcard.com/badge/github.com/johanhenselmans/senbiot)](https://goreportcard.com/report/github.com/johanhenselmans/senbiot)
[![GoDoc](https://godoc.org/github.com/johanhenselmans/senbiot?status.svg)](https://godoc.org/github.com/johanhenselmans/senbiot)
[![codecov](https://codecov.io/gh/johanhenselmans/senbiot/branch/master/graph/badge.svg)](https://codecov.io/gh/johanhenselmans/senbiot)


NB-IOT is a LPWAN technology for bi-directional data traffic between devices and centralized cloud platforms. To configure such a device, a whole slew of AT commands is needed. Senbiot tries to simplify that to a standard set of settings and a simple senbiot "message" command

NOTE: The API is currently unstable

## Installation

To get started please get the Golang toolchains from the [Golang website](https://golang.org/). When you have a working go toolchain you can do:

```
go get github.com/johanhenselmans/senbiot
```

And you are ready to go!

## Included tools

Some simple tools for use with senbiot are included and located in the `cmd` folder of the root of the project.

### Scan serialports (scanserialports)

Commandline tool to scanserialports so as to determine which device to use.

## Contributing

Please read the [Contribution Guidelines](CONTRIBUTING.md). Furthermore: Fork -> Patch -> Push -> Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/johanhenselmans/senbiot/blob/master/LICENSE) file for the full license text.
