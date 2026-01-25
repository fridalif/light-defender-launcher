package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/out
var assets embed.FS

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "My Wails + Next.js App",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("SSH Host: ")
	//host, _ := reader.ReadString('\n')
	//host = strings.TrimSpace(host)
	//
	//fmt.Print("SSH Port (22): ")
	//portInput, _ := reader.ReadString('\n')
	//port := strings.TrimSpace(portInput)
	//if port == "" {
	//	port = "22"
	//}
	//
	//fmt.Print("SSH Username: ")
	//username, _ := reader.ReadString('\n')
	//username = strings.TrimSpace(username)
	//
	//fmt.Print("SSH Password: ")
	//password, _ := term.ReadPassword(int(os.Stdin.Fd()))
	//fmt.Println()
	//var commands []string
	//for {
	//	fmt.Print("Введите команду: ")
	//	cmd, _ := reader.ReadString('\n')
	//	cmd = strings.TrimSpace(cmd)
	//	if cmd == "exit" {
	//		break
	//	}
	//
	//	commands = append(commands, cmd)
	//}
	//sshOperator := sshoperator.NewSSHOperator(username, string(password), host, port, commands)
	//
	//sshOperator.Connect()
	//
	//fmt.Println("\nСессия завершена.")
}
