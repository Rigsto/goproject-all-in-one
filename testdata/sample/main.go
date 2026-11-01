package main

import "fmt"

// greet prints a hello message
func greet(name string) {
	fmt.Println("hello", name)
}

func main() {
	// entry
	greet("world")
}
