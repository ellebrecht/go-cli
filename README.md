
# Geeny Go CLI

## Introduction

Our Go CLI is intended to be the "swiss army knife for Open APIs". We automate and
provide autocompletion when interacting with API's described through Swagger files.

Although this CLI isn't functional with our new platform, it was a great starting
point to automate the creation of CLIs and generate them through their Swagger files
definitions.

## Installation

## Systems Supported:

* Linux `amd64` `386` `arm64` `arm`
* OS X `amd64`
* Windows `amd64` `386`

If you want something else supported, here is
a [list](https://golang.org/doc/install/source#environment) of what can be
supported. Just get in touch with your request.


## Building Source

You will need [golang](https://golang.org/doc/install) installed on your system. We
currently use Go 1.7.4.

Get the code:

```
$ git clone https://github.com/geeny/go-cli.git
```

Add these entries to your rc file, normally `~/.bashrc` if you're using bash

```
export GOPATH=$HOME/go:<YOUR WORKING DIRECTORY>
export GOBIN=$HOME/go/bin
export NETRC=$HOME/.netrc
```

Ensure go is on your `$PATH` with `which go`. If it's not, you may also need to add

```
export PATH=$PATH:/usr/local/bin/go
```

Save changes and reload it with `source ~/.bashrc`

Then do:

* `cd geeny-cli-go`
* `go get ./...` to install dependancies
* `go install geeny` (also see circle.yml for production build flags)
* Run the cli with `$GOBIN/geeny`

> Note: To enable the special debugging option build it with the command: `go build -tags 'debug'`

## Unit Tests

You can run unit tests with `go test <PACKAGE NAME>`

## Make Binaries

Make sure the above steps work, then you can then create cross-platfom binaries with:

* `cd go-cli`
* `chmod a+x makeBinaries.sh`
* `mkdir bin`
* `sh makeBinaries.sh <OPTIONS> ./bin`.

makeBinaries.sh options:
* `-h` displays help
* `-f` force writes over any previous binaries

The directory you provide will contain folders organised by OS/ARCH.

## Configuration

An example of a configuration is provided at `~/.geeny/cli.json`. There you can setup
the services and swagger definitions that should be used when interacting with other
services.

```
{
	"log.enable": false,
	"log.trace": false,
	"log.info": true,
	"log.warn": true,
	"log.error": true,
	"output.spinner": true,
	"config.env": true,
	"update.autocheck": true,
	"swagger.validate": false,
	"net.ca": "",
	"net.proxy": false,
	"net.timeout": 30000000000,
	"output.displayJson": false,
	"output.rawJson": true,
	"connect.url": "https://connect.geeny.io/",
	"connectSwagger.url": "",
	"api.url": "https://api.geeny.io/",
	"apiSwagger.url": "",
	"swagger.defaultContentType": "application/json"
}
```

## Setting up Completions


## IDE

There are many tools to code in golang,
but [Visual Studio Code](https://code.visualstudio.com/) has been useful. The main
plugin you will need is called `Go` by `lukehoban`. The internal plug-in manager will
help you with install this. In general, it will help install a bunch of useful
dependancies, codde anaysis, formatting and so on. This IDE also supports `delve`, a
Go runtime debugger with breakpoints.

# Contributing

## Coding convention

Here's
some
[info](http://stackoverflow.com/questions/22688906/go-naming-conventions-for-const)
on the golang style used throughout the project

## Bugs

Please report any bugs / issues via the
github [issue tracker](https://github.com/geeny/go-cli/issues), with any
information on what you did and what error was displayed

## License

Copyright (C) 2017 Telefónica Germany Next GmbH, Charlottenstrasse 4, 10969 Berlin.

This project is licensed under the terms of
the [Mozilla Public License Version 2.0](LICENSE.md).

Inconsolata font is copyright (C) 2006 The Inconsolata Project Authors. This Font
Software is licensed under the [SIL Open Font License, Version 1.1](OFL.txt).
