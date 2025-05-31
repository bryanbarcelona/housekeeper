package main

import (
	"housekeeper/gui"
)

func main() {
	if err := gui.Run(); err != nil {
		panic(err)
	}
}
