package main

import (
	"fmt"
	"unsafe"
)

func main() {
	num := 5
	numPointer := &num

	// flnum := (*float32)(numPointer)
	flnum := (*float32)(unsafe.Pointer(numPointer))
	fmt.Println(flnum)
}
