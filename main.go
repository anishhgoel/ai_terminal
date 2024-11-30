package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Welcome to custom terminal")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			continue
		}

		input = strings.TrimSpace(input)
		fmt.Println("You entered: ", input)

		if input == ":exit" {
			fmt.Println("Exiting terminal")
			break
		}
		executeCommand(input)

	}
}

func executeCommand(command string) {
	if command == "" {
		return
	}
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("error executing command: %v\n", err)
	}
}
