package sshconnector

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SshConfiguration struct {
	Host     string
	Port     int
	Username string
	Password string
}

type ServiceI interface {
	Connect(sshConfig *SshConfiguration) error
	Execute(command string) (string, error)
	Close() error
	GetSshSession() *ssh.Session
}

type service struct {
	sshClient  *ssh.Client
	sshSession *ssh.Session
}

func NewService() ServiceI {
	return &service{
		sshClient:  nil,
		sshSession: nil,
	}
}

func (s *service) GetSshSession() *ssh.Session {
	return s.sshSession
}

func (s *service) Connect(sshConfig *SshConfiguration) error {
	config := ssh.ClientConfig{
		User: sshConfig.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConfig.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshConfig.Host, sshConfig.Port), &config)
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

func (s *service) Execute(command string) (string, error) {
	if s.sshSession == nil {
		return "", fmt.Errorf("not connected")
	}
	output, err := s.sshSession.CombinedOutput(command)
	return string(output), err
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
