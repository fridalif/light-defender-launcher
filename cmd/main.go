package main

import (
	"bufio"
	"fmt"
	sshoperator "light-defender-launcher/pkg/ssh_operator"
	"os"
	"strings"

	"golang.org/x/term"
)

//func checkSudoPassword(session *ssh.Session) (bool, error) {
//	cmd := "sudo -n true 2>&1"
//
//	defer session.Close()
//
//	output, err := session.CombinedOutput(cmd)
//	if err != nil {
//		return true, nil
//	}
//
//	outputStr := string(output)
//	outputStr = strings.ToLower(outputStr)
//	if strings.Contains(outputStr, "password") ||
//		strings.Contains(outputStr, "sudo") ||
//		strings.Contains(outputStr, "try again") {
//		return true, nil
//	}
//
//	return false, nil
//}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("SSH Host: ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)

	fmt.Print("SSH Port (22): ")
	portInput, _ := reader.ReadString('\n')
	port := strings.TrimSpace(portInput)
	if port == "" {
		port = "22"
	}

	fmt.Print("SSH Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("SSH Password: ")
	password, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	var commands []string
	for {
		fmt.Print("Введите команду: ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "exit" {
			break
		}

		commands = append(commands, cmd)
	}
	sshOperator := sshoperator.NewSSHOperator(username, string(password), host, port, commands)

	sshOperator.Connect()

	fmt.Println("\nСессия завершена.")
}
