package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	mode := "normal"
	fmt.Println("Welcome to custom terminal")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("[%s Mode]> ", strings.Title(mode))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			continue
		}

		input = strings.TrimSpace(input)
		fmt.Println("You entered: ", input)

		if input == ":ai" {
			mode = "ai"
			fmt.Println("Switched to AI mode")
		} else if input == ":normal" {
			mode = "normal"
			fmt.Println("Switched to Normal mode")
		} else if input == ":exit" {
			fmt.Println("Goodbye")
			break
		}

		if mode == "normal" {
			executeCommand(input)
		} else if mode == "ai" {
			processAICommand(input)
		}
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
		fmt.Printf("error executing input: %v\n", err)
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

func suggestSimilarDirectories(inputDir string) {
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
	foundSuggestions := false
	for _, file := range files {
		if file.IsDir() {
			distance := levenshtein.ComputeDistance(file.Name(), inputDir)
			maxLen := float64(max(len(file.Name()), len(inputDir)))
			similarity := 1 - (float64(distance) / maxLen)

			if similarity > 0.5 {
				if !foundSuggestions {
					fmt.Println("Did you mean:")
					foundSuggestions = true
				}
				fmt.Printf(" - %s\n", file.Name())
			}
		}
	}
	if !foundSuggestions {
		fmt.Println("No similar directories found")
	}
}

func processAICommand(command string) {
	if command == "" {
		fmt.Println("Please provide an instruction")
		return
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable")
		return
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		cwd = "unknown"
	}
	dirList := getCurrentDirectories()

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: `You are a coding assistant that translates natural language instructions into shell commands without any explanation. Consider the current working directory and available subdirectories.`,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Current working directory: %s\nAvailable directories: %s\nInstruction: %s", cwd, dirList, command),
		},
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4,
		Messages:    messages,
		MaxTokens:   100,
		Temperature: 0,
	})
	if err != nil {
		fmt.Printf("Failed to get AI suggestion: %v\n", err)
		return
	}
	if len(resp.Choices) == 0 {
		fmt.Println("No command suggested. Please try rephrasing your instruction.")
		return
	}

	suggestedCommand := cleanCommand(resp.Choices[0].Message.Content)
	if suggestedCommand == "" {
		fmt.Println("No command was suggested. Please try rephrasing your instruction.")
		return
	}

	fmt.Printf("Suggested Command: %s\n", suggestedCommand)
	fmt.Print("Execute this command? [y/n]: ")

	reader := bufio.NewReader(os.Stdin)
	userResponse, _ := reader.ReadString('\n')
	userResponse = strings.TrimSpace(strings.ToLower(userResponse))

	if userResponse == "y" || userResponse == "yes" {
		executeCommand(suggestedCommand)
	} else {
		fmt.Println("Command not executed.")
	}
}

// adding this function as some commands generated still had some extra formatting or syntax that is not required or a part of the needed command that is to be executed
func cleanCommand(command string) string {
	command = strings.TrimPrefix(command, "```bash\n")
	command = strings.TrimPrefix(command, "```sh\n")
	command = strings.TrimPrefix(command, "```\n")
	command = strings.TrimSuffix(command, "\n```")
	command = strings.TrimSuffix(command, "```")
	return strings.TrimSpace(command)
}

func getCurrentDirectories() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	files, err := os.ReadDir(currentDir)
	if err != nil {
		return ""
	}

	dirNames := []string{}
	for _, file := range files {
		if file.IsDir() {
			dirNames = append(dirNames, file.Name())
		}
	}
	return strings.Join(dirNames, ", ")

}
