package sshconnector

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type ServiceI interface {
	Connect(host string, port int, username string, password string) error
	Execute(command string) error
	Close() error
}

type service struct {
	commandsChan chan string
	resultsChan  chan string
	sshClient    *ssh.Client
	sshSession   *ssh.Session
	outputChan   chan string
}

func NewService(commandsChan chan string, resultsChan chan string) ServiceI {
	return &service{
		commandsChan: commandsChan,
		resultsChan:  resultsChan,
		sshClient:    nil,
		sshSession:   nil,
	}
}

func (s *service) Connect(host string, port int, username string, password string) error {
	config := ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &config)
	if err != nil {
		return err
	}
	s.sshClient = client
	sshSession, err := client.NewSession()
	if err != nil {
		client.Close()
		return err
	}
	s.sshSession = sshSession
	return nil
}

func (s *service) Execute(command string) error {
	if s.sshSession == nil {
		return fmt.Errorf("not connected")
	}
	output, err := s.sshSession.CombinedOutput(command)
	if err != nil {
		return err
	}
	s.resultsChan <- string(output)
	return nil
}

func (s *service) Close() error {
	if s.sshSession != nil {
		err := s.sshSession.Close()
		s.sshSession = nil
		if err != nil {
			return err
		}
	}
	if s.sshClient != nil {
		err := s.sshClient.Close()
		s.sshClient = nil
		if err != nil {
			return err
		}
	}
	return nil
}
