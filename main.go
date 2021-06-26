package main

import (
	"root.challenge/eventprocessor"
	"root.challenge/input"
)

func main() {
	eventprocessor.New(input.StartReading()).Process()
}
