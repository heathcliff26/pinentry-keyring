package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func createPinentryCMD() *exec.Cmd {
	pinentryPath := os.Getenv("PINENTRY")
	if pinentryPath == "" {
		pinentryPath = "pinentry"
	}
	// #nosec G204: Path is controlled by user environment
	return exec.Command(pinentryPath, os.Args[1:]...)

}

func main() {
	cmd := createPinentryCMD()
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer cmdStdin.Close()
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	defer cmdStdout.Close()
	cmdOut := bufio.NewScanner(cmdStdout)

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	cmdOut.Scan()
	fmt.Println(cmdOut.Text())

	_, err = cmdStdin.Write([]byte("OPTION allow-external-password-cache\n"))
	if err != nil {
		panic(err)
	}

	cmdOut.Scan()
	if cmdOut.Text() != "OK" {
		panic(fmt.Sprintf("Failed to allow external password cache: %s\n", cmdOut.Text()))
	}

	go func() {
		for cmdOut.Scan() {
			fmt.Println(cmdOut.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			// GPG sends "SETKEYINFO --clear" to clear the key info, but we need it to save the key
			if line == "SETKEYINFO --clear" {
				line = "SETKEYINFO pinenty-keyring-default-key"
			}
			_, err := cmdStdin.Write([]byte(line + "\n"))
			if err != nil {
				panic(err)
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}
