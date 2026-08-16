package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// #nosec G703: Test binary, input is controlled by caller.
	f, err := os.Create(os.Getenv("PINENTRY_MOCK_OUTPUT"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	response := "OK Pleased to meet you\n"
	_, _ = f.WriteString(response)
	_, _ = fmt.Print(response)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		response := ""
		exit := false
		switch line {
		case "GETPIN":
			response = "S PASSWORD_FROM_CACHE\nD 1234\n"
		case "BYE":
			response = "OK closing connection\n"
			exit = true
		default:
			response = "OK\n"
		}

		_, _ = f.WriteString(line + "\n" + response)
		_, _ = fmt.Print(response)
		if exit {
			break
		}
	}
}
