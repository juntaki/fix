# fix - Golden Files Testing library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/juntaki/fix.svg)](https://pkg.go.dev/github.com/juntaki/fix)

This library, like Ruby's [vcr](https://github.com/vcr/vcr), serializes the results once executed and saves them as a file. You can easily write tests to verify that the code outputs the same results as in the past, even after modifying the code.

~~~
go get github.com/juntaki/fix
~~~

## How to use

Append code like below to your tests. and Run test twice.

~~~
err := fix.Fix(&output) // output is dumped to a file.
if err != nil {
  t.Fatal(err)
}
~~~

First test will fail, because juntaki/fix writes serialized data to file.
From the second time, the test will pass, if output is the same as first output.
