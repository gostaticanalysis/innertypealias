[![pkg.go.dev][gopkg-badge]][gopkg]

`innertypealias` finds find a type which is an alias for exported same package's type.

```go
package a

import "io"

type T int
type t int

type A = T         // want "A is an alias for T but it is exported type"
type B = t         // OK
type C = io.Writer // OK

func _() {
	type D = T // OK
}

type E T   // OK
type F t   // OK
type g = t // OK

type H = T // OK - it is used as an embedded field
type _ struct{ H }

type I = T // OK - it is used as an embedded field
func _() {
	type _ struct{ I }
}

type _ = T // OK
```

`fixinnertypealias` command check and replace a type alias to a defined type.

```sh
$ go install github.com/gostaticanalysis/innertypealias/cmd/fixinnertypealias@latest
$ fixinnertypealias ./...
```
<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gostaticanalysis/innertypealias
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gostaticanalysis/innertypealias?status.svg
