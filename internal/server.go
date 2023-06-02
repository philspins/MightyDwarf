package internal

import (
	"fmt"
	"io"
	"net"
	"strings"
	"unicode"
)

type Server struct {
	ipAddress           string
	port                int
	TelescopeController *TelescopeController
}

func NewLX200Server(port int) *Server {
	return &Server{
		ipAddress:           "0.0.0.0",
		port:                port,
		TelescopeController: NewTelescopeController(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ipAddress, s.port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		client, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handleClient(client)
	}
}

func (s *Server) handleClient(client net.Conn) {
	defer client.Close()

	reader := strings.NewReader("")
	writer := client

	buffer := make([]byte, 1)
	command := strings.Builder{}

	for {
		_, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		receivedChar := rune(buffer[0])

		if unicode.IsSpace(receivedChar) {
			continue
		}

		if receivedChar == '#' {
			// Command terminator received, process the command.
			response := s.TelescopeController.HandleCommand(command.String())
			if command.String() != ":GR" && command.String() != ":GD" {
				fmt.Println(command.String())
			}

			if response != "" {
				writer.Write([]byte(response))
			}

			command.Reset()
		} else {
			command.WriteRune(receivedChar)
		}
	}
}
