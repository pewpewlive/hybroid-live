//go:build js && wasm

package main

import (
	_ "hybroid/wasm"
)

func main() {
	// Prevent the program from exiting so the JS functions remain registered
	select {}
}
