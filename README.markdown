Erlang External Term Format for Go
==================================

[![Build Status](https://app.travis-ci.com/okeuday/erlang_go.svg?branch=master)](https://app.travis-ci.com/okeuday/erlang_go) [![Go Report Card](https://goreportcard.com/badge/github.com/okeuday/erlang_go?maxAge=3600)](https://goreportcard.com/report/github.com/okeuday/erlang_go)

Provides all encoding and decoding for the Erlang External Term Format
(as defined at [https://erlang.org/doc/apps/erts/erl_ext_dist.html](https://erlang.org/doc/apps/erts/erl_ext_dist.html))
in a single Go package.

(For `go` (version < 1.11) command-line use you can use the prefix
 `GOPATH=`pwd` GOBIN=$$GOPATH/bin` to avoid additional shell setup)
(For `go` (version > 1.11) command-line use you can use the prefix
 `GO111MODULE=auto` to avoid additional shell setup)

Build
-----

    go build ./...

Test
----

    go test ./...

Author
------

Michael Truog (mjtruog at protonmail dot com)

License
-------

MIT License

