# Packman

A rudimentary package indexer that keeps track of package dependencies. Note that versioning of packages is not part of this initial requirements. See the _Requirements_ section below for detailed requirements for how this indexer works.

[![Go Report Card](https://goreportcard.com/badge/github.com/jboursiquot/packman)](https://goreportcard.com/report/github.com/jboursiquot/packman) [![Build Status](https://travis-ci.org/jboursiquot/packman.svg?branch=master)](https://travis-ci.org/jboursiquot/packman)

## Install

`Packman` is go-gettable so simply enter this in your terminal:

```plain
$ go get github.com/jboursiquot/packman/cmd/packman
```

## Usage

To start `Packman`:

```plain
$ packman
2016/07/30 17:39:30 Starting server...
```

Packman listens on port `8080` so your OS may require you to grant the ability for the executable to use that port.

## Logging

`Packman` will emit logs to `stdout` as requests are sent to it and as it responds.

Sample logging output:

```plain
2016/07/30 17:56:12 REQ=INDEX|hubflow|git, RES=OK
2016/07/30 17:56:12 REQ=REMOVES|a|b, RES=ERROR
2016/07/30 17:56:12 REQ=INDEX|unac|autoconf,automake,gettext,libtool,xz, RES=OK
2016/07/30 17:56:12 REQ=INDEX|emacs☃elisp, RES=ERROR
2016/07/30 17:56:12 REQ=INDEX|memo|, RES=OK
2016/07/30 17:56:12 REQ=INDEX|ortp|, RES=OK
2016/07/30 17:56:12 REQ=INDEX|openhmd|autoconf,automake,cmake,hidapi,libtool,pkg-config,sphinx-doc,xz, RES=OK
2016/07/30 17:56:12 REQ=INDEX|dex|, RES=OK
2016/07/30 17:56:12 REQ=INDEX|libb2|, RES=OK
2016/07/30 17:56:12 REQ=INDEX|jcal|autoconf,automake,libtool,xz, RES=OK
```

## Tests

Nothing special here. The standard `go` tools work:

```plain
$ go test -race -v
=== RUN   TestHandlerHandlesInvalidMessages
--- PASS: TestHandlerHandlesInvalidMessages (0.00s)
=== RUN   TestHandlerTranslatesMessageToCommand
--- PASS: TestHandlerTranslatesMessageToCommand (0.00s)
=== RUN   TestHandlerProcessesCommandAccurately
--- PASS: TestHandlerProcessesCommandAccurately (0.00s)
=== RUN   TestIndexerIndexesPackageWithoutDeps
--- PASS: TestIndexerIndexesPackageWithoutDeps (0.00s)
=== RUN   TestIndexerIndexesPackageWithKnownDeps
--- PASS: TestIndexerIndexesPackageWithKnownDeps (0.00s)
=== RUN   TestIndexerFailsToIndexPackageWithUnknownDeps
--- PASS: TestIndexerFailsToIndexPackageWithUnknownDeps (0.00s)
=== RUN   TestIndexerRemovesPackageWithoutDeps
--- PASS: TestIndexerRemovesPackageWithoutDeps (0.00s)
=== RUN   TestIndexerFailsToRemovePackageWithDeps
--- PASS: TestIndexerFailsToRemovePackageWithDeps (0.00s)
=== RUN   TestIndexerFindsIndexedPackage
--- PASS: TestIndexerFindsIndexedPackage (0.00s)
=== RUN   TestIndexerFailsToFindUnindexedPackage
--- PASS: TestIndexerFailsToFindUnindexedPackage (0.00s)
PASS
ok  	github.com/jboursiquot/packman	1.019s
```

## Requirements

Clients will connect to the server and inform which packages should be indexed, and which dependencies they might have on other packages. We want to keep our index consistent, so your server must not index any package until all of its dependencies have been indexed first. The server should also not remove a package if any other packages depend on it.

The server will open a TCP socket on port 8080. It must accept connections from multiple clients at the same time, all trying to add and remove items to the index concurrently. Clients are independent of each other, and it is expected that they may send repeated or contradicting messages. New clients can connect and disconnect at any moment, and sometimes clients can behave badly and try to send broken messages.

Messages from clients follow this pattern:

```plain
<command>|<package>|<dependencies>\n
```

Where:

* `<command>` is mandatory, and is either `INDEX`, `REMOVE`, or `QUERY`
* `<package>` is mandatory, the name of the package referred to by the command, e.g. `mysql`, `openssl`, `pkg-config`, `postgresql`, etc.
* `<dependencies>` is optional, and if present it will be a comma-delimited list of packages that need to be present before `<package>` is installed. e.g. `cmake`,`sphinx-doc`,`xz`
* The message always ends with the character `\n`

Sample messages:

```plain
INDEX|cloog|gmp,isl,pkg-config\n
INDEX|ceylon|\n
REMOVE|cloog|\n
QUERY|cloog|\n
```

For each message sent, the client will wait for a response code from the server. Possible response codes are `OK\n`, `FAIL\n`, or `ERROR\n`. After receiving the response code, the client can send more messages.

The response code returned should be as follows:

* For `INDEX` commands, the server returns `OK\n` if the package can be indexed. It returns `FAIL\n` if the package cannot be indexed because some of its dependencies aren’t indexed yet and need to be installed first. If a package already exists, then its list of dependencies is updated to the one provided with the latest command.

* For `REMOVE` commands, the server returns `OK\n` if the package could be removed from the index. It returns `FAIL\n` if the package could not be removed from the index because some other indexed package depends on it. It returns `OK\n` if the package wasn’t indexed.

* For `QUERY` commands, the server returns `OK\n` if the package is indexed. It returns `FAIL\n` if the package isn’t indexed.

* If the server doesn’t recognize the command or if there’s any problem with the message sent by the client it should return `ERROR\n`.
