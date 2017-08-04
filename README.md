# go-tarantool [![GoDoc](https://godoc.org/github.com/viciious/go-tarantool?status.svg)](https://godoc.org/github.com/viciious/go-tarantool) [![Build Status](https://travis-ci.org/viciious/go-tarantool.svg?branch=master)](https://travis-ci.org/viciious/go-tarantool)

<a href="http://tarantool.org">
	<img src="https://avatars2.githubusercontent.com/u/2344919?v=2&s=250" align="right">
</a>

The `go-tarantool` package has everything necessary for interfacing with
[Tarantool 1.6+](http://tarantool.org/).

The advantage of integrating Go with Tarantool, which is an application server
plus a DBMS, is that Go programmers can handle databases and perform on-the-fly
recompilations of embedded Lua routines, just as in C, with responses that are
faster than other packages according to public benchmarks.

## Table of contents

* [Installation](#installation)
* [Hello World](#hello-world)
* [API reference](#api-reference)
* [Walking through the example](#walking-through-example-in-go)
* [Help](#help)
* [Alternative way to connect](#alternative-way-to-connect)

## Installation

We assume that you have Tarantool version 1.6 or 1.7, and a modern Linux or BSD
or Mac OS operating system.

You will need a current version of `go`, version 1.8 or later (use
`go version` to check the version number).

If your `go` version is older than 1.8, or if `go` is not installed,
download the latest tarball from [golang.org](https://golang.org/dl/) and say:

```bash
sudo tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz
sudo chmod -R a+rwx /usr/local/go
```

Make sure `go` and `go-tarantool` are on your path. For example:

```
export PATH=$PATH:/usr/local/go/bin
export GOPATH="/usr/local/go/go-tarantool"
```

The `go-tarantool` package is in the 
[viciious/go-tarantool](https://github.com/viciious/go-tarantool) repository.
To download and install, say:

```
go get github.com/viciious/go-tarantool
```

This should bring source and binary files into subdirectories of `/usr/local/go`,
making it possible to access by adding `github.com/viciious/go-tarantool` in
the `import {...}` section at the start of any Go program.

<h2>Hello World</h2>

Here is a very short example Go program which tries to connect to a Tarantool server.

```go
package main

import (
    "context"
    "fmt"
    "github.com/viciious/go-tarantool"
)

func main() {
    opts := tarantool.Options{User: "guest"}
    conn, err := tarantool.Connect("127.0.0.1:3301", &opts)
    if err != nil {
        fmt.Println("Connection refused: %s", err.Error())
    }

    query := &tarantool.Insert{Space: "examples", Tuple: []interface{}{uint64(99999), "BB"}}
    resp := conn.Exec(context.Background(), query)

    if resp.Error != nil {
        fmt.Println("Insert failed", resp.Error)
    } else {
        fmt.Println(fmt.Sprintf("Insert succeeded: %#v", resp.Data))
    }

    conn.Close()
}
```

Cut and paste this example into a file named example.go.

Start a Tarantool server on localhost, and make sure it is listening
on port 3301. Set up a space named `examples` exactly as described
in the Tarantool manual's Connectors section
[Tarantool manual's Connectors section](https://tarantool.org/doc/1.7/book/connectors/index.html#index-connector-setting).

Again, make sure PATH and GOPATH point to the right places. Then build and run `example.go`:

```
go build example.go
./example
```

You should see: messages saying "Insert failed" or "Insert succeeded".

If that is what you see, then you have successfully installed `go-tarantool` and
successfully executed a program that connected to a Tarantool server and
manipulated the contents of a Tarantool database.


## Walking through the example

We can now have a closer look at the `example.go` program and make some observations
about what it does.

**Observation 1:** the line "`github.com/viciious/go-tarantool`" in the
`import(...)` section brings in all Tarantool-related functions and structures.
It is common to bring in [context](https://golang.org/pkg/context/)
and [fmt](https://golang.org/pkg/fmt/) as well.

**Observation 2:** the line beginning with "`Opts :=`" sets up the options for
`Connect()`. In this example, there is only one thing in the structure, a user
name. The structure can also contain:

* `ConnectTimeout` 
* `QueryTimeout`    (the default maximum number of milliseconds to wait before giving up - can be overriden on per-query basis),
* `DefaultSpace`    (the name of default Tarantool space)
* `Password`        (user's password)
* `UUID`            (used for replication)
* `ReplicaSetUUID`  (used for replication)

**Observation 3:** the line containing "`tarantool.Connect`" is one way
to begin a session. There are two parameters:

* a string with `host:port` format (or "/path/to/tarantool.socket"), and
* the option structure that was set up earlier.

There is an alternative way to connect, we will describe it later.

**Observation 4:** the `err` structure will be `nil` if there is no error,
otherwise it will have a description which can be retrieved with `err.Error()`.

**Observation 5:** the `conn.exec`request, like many requests, is preceded by
"`conn.`" which is the name of the object that was returned by `Connect()`.
In this case, for Insert, there are two parameters:

* a space name (it could just as easily have been a space number), and
* a tuple.

All the requests described in the Tarantool manual can be expressed in
a similar way within connect.Exec(), with the format "&name-of-request{arguments}".
For example: `&ping{}`. For a long example:

```go
    data, err := conn.Exec(context.Background(), &Update{
        Space: "tester",
        Index: "primary",
        Key:   1,
        Set: []Operator{
            &OpAdd{
                Field:    2,
                Argument: 17,
            },
            &OpAssign{
                Field:    1,
                Argument: "Hello World",
            },
        },
    })
```

<h2>API reference</h2>

Read the [Tarantool manual](http://tarantool.org/doc.html) to find descriptions
of terms like "connect", "space", "index", and the requests for creating and
manipulating database objects or Lua functions.

The source files for the requests library are:
* [connection.go](https://github.com/viciious/go-tarantool/blob/master/connector.go)
  for the `Connect()` function plus functions related to connecting, and
* [insert_test.go](https://github.com/viciious/go-tarantool/blob/master/insert_test.go)
  for an example of a data-manipulation function used in tests.

See comments in these files for syntax details:
```
call.go
delete.go
eval.go
insert.go
iterator.go
join.go
operator.go
pack.go
update.go
upsert.go
```

The supported requests have parameters and results equivalent to requests in the
Tarantool manual. Browsing through the other *.go programs in the package will
show how the packagers have paid attention to some of the more advanced features
of Tarantool, such as vclock and replication.

## Alternative way to connect

Here we show a variation of example.go, where the connect is
done a different way.

```go

package main

import (
    "context"
    "fmt"
    "github.com/viciious/go-tarantool"
)

func main() {
    opts := tarantool.Options{User: "guest"}
    conn, err := tarantool.Connect("127.0.0.1:3301", &opts)
    if err != nil {
        fmt.Println("Connection refused: %s", err.Error())
    }

    query := &tarantool.Insert{Space: "examples", Tuple: []interface{}{uint64(99999), "BB"}}
    resp := conn.Exec(context.Background(), query)

    if resp.Error != nil {
        fmt.Println("Insert failed", resp.Error)
    } else {
        fmt.Println(fmt.Sprintf("Insert succeeded: %#v", resp.Data))
    }

    conn.Close()
}
```

In this variation, tarantool.New returns a Connector instance,
which is a goroutine-safe object that can transparently handle
reconnects.

## Help

To contact `go-tarantool` developers on any problems, create an issue at
[viciious/go-tarantool](http://github.com/viciious/go-tarantool/issues).

The developers of the [Tarantool server](http://github.com/tarantool/tarantool)
will also be happy to provide advice or receive feedback.
