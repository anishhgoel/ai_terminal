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

	if args[0] == "cd" {
		changeDirectory(args)
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("error executing command: %v\n", err)
	}
}

// cmd.Run() executes commands in child processes, which are isolated from the main program so need to make a function for cd to change directories
func changeDirectory(args []string) {
	var dir string
	if len(args) < 2 {
		dir = os.Getenv("HOME")
		if dir == "" {
			fmt.Println("HOME environment variable is not set.")
			return
		}
	} else {
		dir = args[1]
	}
	dir = os.ExpandEnv(dir) //expandung the environment variables (e.g. $HOME)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist %s\n", dir)
		suggestSimilarDirectories(dir)
		return
	}
	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("Error changing directory: %v \n", err)
		return
	}
	fmt.Printf("Changed directory to %s", dir)

}

func suggestSimilarDirectories(dir string) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting the current directory %v\n", err)
		return
	}

	files, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Printf("Error getting the files of current direcotry %v\n", err)
		return
	}
	fmt.Println("Did you mean: ")
	for _, file := range files {
		if file.IsDir() {
			fmt.Printf(" - %s\n", file.Name())
		}
	}

}
