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
