package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
)

var (
	outputBuf = io.Writer(os.Stdout)
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) checkSudoPassword(session *ssh.Session) (bool, error) {
	cmd := "sudo -n true 2>&1"

	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return true, nil
	}

	outputStr := string(output)
	outputStr = strings.ToLower(outputStr)
	if strings.Contains(outputStr, "password") ||
		strings.Contains(outputStr, "sudo") ||
		strings.Contains(outputStr, "try again") {
		return true, nil
	}

	return false, nil
}

func (a *App) ConnectAndExecuteCommands(username string, password string, host string, port string, commands []string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(password)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Подключаемся к %s...\n", address))
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("Ошибка подключения: %v", err)
	}
	defer client.Close()
	runtime.EventsEmit(a.ctx, "log", "Подключение установлено!") //fmt.Println("Подключение успешно!")

	testsession, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Ошибка создания сессии: %v", err)
	}
	var needRootPass = false
	needRootPass, err = a.checkSudoPassword(testsession)
	if err != nil {
		return fmt.Errorf("Ошибка проверки sudo: %v", err)
	}
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Ошибка создания сессии: %v", err)
	}

	//fd := int(os.Stdin.Fd())
	//oldState, err := term.MakeRaw(fd)
	//if err != nil {
	//	return fmt.Errorf("Ошибка получения состояния терминала: %v", err)
	//}
	//defer term.Restore(fd, oldState)

	termWidth, termHeight := 0, 0
	if termWidth == 0 {
		termWidth = 80
	}
	if termHeight == 0 {
		termHeight = 24
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", termHeight, termWidth, modes)
	if err != nil {
		return fmt.Errorf("Ошибка запроса pty: %v", err)
	}

	session.Stdout = outputBuf
	session.Stderr = outputBuf
	stdin, _ := session.StdinPipe()

	err = session.Shell()
	if err != nil {
		return fmt.Errorf("Ошибка запуска shell: %v", err)
	}

	runtime.EventsEmit(a.ctx, "log", "Выполняем sudo su root...") //fmt.Println("\nВыполняем sudo su root...")

	stdin.Write([]byte("sudo su root\n"))
	time.Sleep(1 * time.Second)
	if needRootPass {
		stdin.Write([]byte(string(password) + "\n"))
	}
	time.Sleep(2 * time.Second)

	for _, command := range commands {
		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Выполняем команду: %s\n", command))
		stdin.Write([]byte(command + "\n"))
		time.Sleep(2 * time.Second)
	}

	err = session.Close()
	if err != nil {
		return nil
	}
	return nil
}

type UserInput struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfigurationID string `json:"configuration_id"`
}

func LoadFileFromDashboard(dashboardLogin string, dashboardPassword string, configId string) ([]byte, error) {
	url := "https://bot.light-defender.ru/load_config"

	userInput := UserInput{
		Username:        dashboardLogin,
		Password:        dashboardPassword,
		ConfigurationID: configId,
	}
	inputBytes, err := json.Marshal(userInput)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(inputBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Учётные данные не верны.")
	}

	return body, nil
}

func (a *App) LoadConfig(sshUsername string, sshPassword string, host string, port string, dashboardLogin string, dashboardPassword string, configId string, fileb64 string) []string {
	if ((dashboardLogin == "") || (dashboardPassword == "") || (configId == "")) && len(fileb64) == 0 {
		return []string{"false", "Поля не заполнены"}
	}
	if len(fileb64) == 0 {
		file, err := LoadFileFromDashboard(
			dashboardLogin,
			dashboardPassword,
			configId,
		)
		if err != nil {
			return []string{"false", err.Error()}
		}
		fileb64 = base64.StdEncoding.EncodeToString(file)
	}
	commands := []string{
		`[[ -f /etc/systemd/system/ldclient.service ]] && cd "$(grep WorkingDirectory /etc/systemd/system/ldclient.service | cut -d= -f2 | xargs)" 2>/dev/null && pwd && ` +
			fmt.Sprintf(`echo '%s' | base64 -d > ./etc/config.bin`, fileb64) +
			` && sudo systemctl restart ldclient || echo "Файл не найден или ошибка"`,
	}

	err := a.ConnectAndExecuteCommands(sshUsername, sshPassword, host, port, commands)

	if err != nil {
		return []string{"false", err.Error()}
	}

	return []string{"true", ""}
}

func (a *App) InstallLightDefender(sshUsername string, sshPassword string, host string, port string) []string {
	commands := []string{
		`[[ ! -f /etc/systemd/system/ldclient.service ]] && wget -q -O - https://dashboard.light-defender.ru/auto_load_ld.sh | sudo bash || echo "Файл не найден или ошибка"`,
	}

	err := a.ConnectAndExecuteCommands(sshUsername, sshPassword, host, port, commands)

	if err != nil {
		return []string{"false", err.Error()}
	}

	return []string{"true", ""}
}

func (a *App) UpdateLightDefender(sshUsername string, sshPassword string, host string, port string) []string {
	commands := []string{
		`[[ -f /etc/systemd/system/ldclient.service ]] && cd "$(grep WorkingDirectory /etc/systemd/system/ldclient.service | cut -d= -f2 | xargs)" 2>/dev/null && pwd && sudo systemctl stop ldclient && sudo mv ldclient.bin ldclient.bin.save && sudo ./ldclient.bin.save -uni && sudo rm -rf ldclient.bin.save && sudo systemctl start ldclient || echo "Файл не найден или ошибка"`,
	}

	err := a.ConnectAndExecuteCommands(sshUsername, sshPassword, host, port, commands)

	if err != nil {
		return []string{"false", err.Error()}
	}

	return []string{"true", ""}
}
