package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	stateFileName := "test"

	state := NewState("data", "csv")
	fmt.Println(state.String() + "\n")

	err := state.LoadEntries(stateFileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(state.String() + "\n")

	entry := NewEntry(
		"cmpny",
		"dev",
		"applied",
		"linkeding",
		time.Now(),
		false,
	)
	entry.AddToState(state)
	state.SaveEntries(stateFileName)
	fmt.Println(state.String() + "\n")
}

