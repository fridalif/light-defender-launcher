package main

import (
	"fmt"
	sshconnector "light-defender-launcher/pkg/ssh_connector"
	"os"

	"golang.org/x/term"
)

func main() {
	var config sshconnector.SshConfiguration
	fmt.Print("Имя пользователя: ")
	fmt.Scanln(&config.Username)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("\nError reading password: %v\n", err)
		return
	}
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
	for {
		fmt.Print("Введите команду: ")
		var command string
		fmt.Scanln(&command)
		result, err := sshService.Execute(command)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(result)
	}
}
