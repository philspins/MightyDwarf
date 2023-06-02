package internal

import (
	"fmt"
	"io"
	"log"

	"github.com/tarm/serial"
)

type SerialPortListener struct {
	serialPort        *serial.Port
	telescopeControl  *TelescopeController
	readBuffer        []byte
	readBufferIndex   int
	readBufferMaxSize int
}

func NewSerialPortListener(comPortName string, telescopeControl *TelescopeController) (*SerialPortListener, error) {
	c := &serial.Config{
		Name:     comPortName,
		Baud:     9600,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
		Size:     8,
		// Handshake: serial.HandshakeNone,
	}

	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	listener := &SerialPortListener{
		serialPort:        s,
		telescopeControl:  telescopeControl,
		readBuffer:        make([]byte, 128),
		readBufferIndex:   0,
		readBufferMaxSize: 128,
	}

	go listener.startReading()

	return listener, nil
}

func (l *SerialPortListener) startReading() {
	for {
		n, err := l.serialPort.Read(l.readBuffer[l.readBufferIndex:l.readBufferMaxSize])
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from serial port:", err)
			}
			return
		}

		l.readBufferIndex += n

		if l.readBufferIndex >= l.readBufferMaxSize {
			log.Println("Serial port read buffer is full")
			return
		}

		l.processReceivedData()
	}
}

func (l *SerialPortListener) processReceivedData() {
	startIndex := 0
	for i := 0; i < l.readBufferIndex; i++ {
		if l.readBuffer[i] == '#' {
			command := string(l.readBuffer[startIndex:i])
			if command != ":GR" && command != ":GD" {
				fmt.Println(command)
			}

			l.telescopeControl.HandleCommand(command)
			startIndex = i + 1
		}
	}

	// Shift the remaining data in the buffer to the beginning
	copy(l.readBuffer, l.readBuffer[startIndex:l.readBufferIndex])
	l.readBufferIndex -= startIndex
}

func (l *SerialPortListener) Start() {
	// Nothing to do here since reading is already started in a separate goroutine
}

func (l *SerialPortListener) Stop() {
	err := l.serialPort.Close()
	if err != nil {
		log.Println("Error closing serial port:", err)
	}
}
