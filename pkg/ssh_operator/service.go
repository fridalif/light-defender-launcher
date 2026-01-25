package sshoperator

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHOperatorI interface {
	Connect()
}

type SSHOperator struct {
	username string
	password string
	host     string
	port     string
	commands []string
}

func NewSSHOperator(username, password, host, port string, commands []string) SSHOperatorI {
	return &SSHOperator{
		username: username,
		password: password,
		host:     host,
		port:     port,
		commands: commands,
	}
}

func (s *SSHOperator) checkSudoPassword(session *ssh.Session) (bool, error) {
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

func (s *SSHOperator) Connect() {
	address := fmt.Sprintf("%s:%s", s.host, s.port)
	config := &ssh.ClientConfig{
		User: s.username,
		Auth: []ssh.AuthMethod{
			ssh.Password(string(s.password)),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	fmt.Printf("Подключаемся к %s...\n", address)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer client.Close()
	fmt.Println("Подключение успешно!")

	testsession, err := client.NewSession()
	if err != nil {
		log.Fatalf("Ошибка создания сессии: %v", err)
	}
	var needRootPass = false
	needRootPass, err = s.checkSudoPassword(testsession)
	if err != nil {
		log.Fatal(err)
	}
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Ошибка создания сессии: %v", err)
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(fd, oldState)

	termWidth, termHeight, _ := term.GetSize(fd)
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
		log.Fatalf("Ошибка запроса pty: %v", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	stdin, _ := session.StdinPipe()

	err = session.Shell()
	if err != nil {
		log.Fatalf("Ошибка запуска shell: %v", err)
	}

	fmt.Println("\nВыполняем sudo su root...")

	stdin.Write([]byte("sudo su root\n"))
	time.Sleep(1 * time.Second)
	if needRootPass {
		stdin.Write([]byte(string(s.password) + "\n"))
	}
	time.Sleep(2 * time.Second)

	for _, command := range s.commands {
		stdin.Write([]byte(command + "\n"))
		time.Sleep(2 * time.Second)
	}

	session.Close()
}
