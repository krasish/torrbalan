package command

import (
	"bufio"
	"fmt"
	"os"
)

type Interpreter struct {
}

func (i Interpreter) GetUsername() (username string) {
	var (
		err    error
		reader = bufio.NewReader(os.Stdin)
	)

	for {
		fmt.Print("Please enter a username: ")
		username, err = reader.ReadString('\n')
		if err == nil {
			break
		}
		fmt.Println("That didn't work. Try again!")
	}

	return
}
