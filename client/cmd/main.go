package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("could not dial: %v", err)
	}
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			bytes := make([]byte, 100)
			n, err2 := conn.Read(bytes)
			if err2 != nil {
				log.Fatalf("while reading: %v read %d", err, n)
			}
			fmt.Println(string(bytes))
		}
	}()

	for {
		fmt.Print("> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("could not read: %v", err)
		continue
	}
		if _, err := conn.Write([]byte(text)); err != nil {
			fmt.Printf("could not write: %v", err)
		}

	}
}
