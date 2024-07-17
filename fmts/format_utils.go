package fmts

import "fmt"

func Panicf(format string, v ...any) {
	panic(fmt.Sprintf(format, v...))
}
