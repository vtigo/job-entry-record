package main

import "fmt"

func main() {
	state := NewState("data")
	state.LoadEntries("test")
	fmt.Println(state.String())
}
