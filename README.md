[![pkg.go.dev][gopkg-badge]][gopkg]

`innertypealias` finds find a type which is an alias for exported same package's type.

```go
package a

import "io"

type T int
type t int

type A = T         // want "A is a alias for T but it is exported type"
type B = t         // OK
type C = io.Writer // OK

func f() {
	type D = T // OK
}

type E T   // OK
type F t   // OK
type g = t // OK
```

`fixinnertypealias` command check and replace a type alias to a defined type.

```sh
$ go install github.com/gostaticanalysis/innertypealias/cmd/fixinnertypealias
$ fixinnertypealias ./...
```
<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gostaticanalysis/innertypealias
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gostaticanalysis/innertypealias?status.svg
