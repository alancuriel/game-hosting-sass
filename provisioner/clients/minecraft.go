package clients

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type McServer struct {
	ip        string
	sshConfig *ssh.ClientConfig
}

func MCServer(ip, sshUser, sshPsswd string) *McServer {
	return &McServer{
		ip: ip,
		sshConfig: &ssh.ClientConfig{
			User: sshUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(sshPsswd),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         10 * time.Second,
		},
	}
}

func (s *McServer) Announce(message string) error {
	runnable := func(session *ssh.Session) error {
		// Escape any quotes in the message to prevent command injection
		escapedMessage := strings.Replace(message, "'", "\\'", -1)
		escapedMessage = strings.Replace(escapedMessage, "\"", "\\\"", -1)

		// Execute the minecraft server command
		cmd := fmt.Sprintf("su - mcserver -c '/home/mcserver/mcserver send \"say %s\"'", escapedMessage)
		err := session.Run(cmd)
		if err != nil {
			return fmt.Errorf("failed to execute send command: %v", err)
		}

		return nil
	}

	return s.createSessionAndRun(runnable)
}

func (s *McServer) Stop() error {
	runnable := func(session *ssh.Session) error {
		err := session.Run("su - mcserver -c '/home/mcserver/mcserver stop")

		if err != nil {
			return fmt.Errorf("failed to execute stop command: %v", err)
		}

		return nil
	}

	return s.createSessionAndRun(runnable)
}

func (s *McServer) Start() error {
	runnable := func(session *ssh.Session) error {
		err := session.Run("su - mcserver -c '/home/mcserver/mcserver start")

		if err != nil {
			return fmt.Errorf("failed to execute start command: %v", err)
		}

		return nil
	}

	return s.createSessionAndRun(runnable)
}

func (s *McServer) createSessionAndRun(runner func(session *ssh.Session) error) error {
	// Connect to the server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", s.ip), s.sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create session
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	err = runner(session)

	if err != nil {
		return err
	}
	return nil
}
