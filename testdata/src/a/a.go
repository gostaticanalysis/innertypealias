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
