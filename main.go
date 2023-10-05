package main

import (
	"./greeting"
	"github.com/fatih/color"
)

func main() {
	color.Red(greeting.Hello())
}
