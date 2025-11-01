package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/s-bose/mox/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Welcome to the interactive prompt for the mox language\n", user.Username)

	repl.Start(os.Stdin, os.Stdout)
}
