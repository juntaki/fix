# fix
[![Build Status](https://travis-ci.org/juntaki/fix.svg?branch=master)](https://travis-ci.org/juntaki/fix) [![GoDoc](https://godoc.org/github.com/juntaki/fix?status.svg)](https://godoc.org/github.com/juntaki/fix) [![Go Report Card](https://goreportcard.com/badge/github.com/juntaki/fix)](https://goreportcard.com/report/github.com/juntaki/fix)    

~~~
go get github.com/juntaki/fix
~~~

## How to use

Append code like below to your tests. and Run test twice.

~~~
err := fix.Fix(&output)
if err != nil {
  t.Fatal(err)
}
~~~

First test will fail, because juntaki/fix writes gob encoded binary to file.
From the second time, the test will pass, if output is the same as first output.
