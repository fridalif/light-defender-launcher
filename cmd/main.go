package main

import (
	"bufio"
	"fmt"
	sshconnector "light-defender-launcher/pkg/ssh_connector"
	"os"

	"golang.org/x/term"
)

func main() {
	var config sshconnector.SshConfiguration
	fmt.Print("Имя пользователя: ")
	fmt.Scanln(&config.Username)
	fmt.Print("Пароль: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("\nError reading password: %v\n", err)
		return
	}
	fmt.Println()
	config.Password = string(password)
	fmt.Print("IP: ")
	fmt.Scanln(&config.Host)
	fmt.Print("Порт: ")
	fmt.Scanln(&config.Port)
	sshService := sshconnector.NewService()
	err = sshService.Connect(&config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sshService.Close()
	stdOut, err := sshService.GetSshSession().StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	stdIn, err := sshService.GetSshSession().StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := sshService.GetSshSession().Shell(); err != nil {
		fmt.Println(err.Error())
		return
	}
	go func() {
		scanner := bufio.NewScanner(stdOut)
		for scanner.Scan() {
			fmt.Println("Remote Output:", scanner.Text())
		}
	}()

	for {
		fmt.Print("Команда(или exit): ")
		var command string
		fmt.Scanln(&command)
		if command == "exit" {
			break
		}
		_, err := stdIn.Write([]byte(command + "\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	stdIn.Close()
}
