package main

import (
	"tmail/internal"
)

func main() {
	if err := internal.NewApp().Run(); err != nil {
		panic(err)
	}
}
