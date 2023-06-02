package internal

import (
	"errors"
	"fmt"
	"time"
)

type dwarf2Client struct {
	webSocketClient *websock
}

func (c *dwarf2Client) init(lon, lat float64, date time.Time, path string) error {
	c.webSocketClient = &websock{Config: Config{IpAddress: "10.0.0.8"}}
	err := c.webSocketClient.Connect()
	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
		return err
	}

	err = c.webSocketClient.UpdateDateTime(date)
	if err != nil {
		fmt.Println("Error updating date and time:", err)
		return err
	}

	// Perform the "correction" request
	err = c.webSocketClient.SendCorrection(lon, lat, date, path)
	if err != nil {
		fmt.Println("Error sending correction:", err)
		return err
	}

	fmt.Println("Correction request succeeded.")

	return nil
}

func (c *dwarf2Client) Goto(lon, lat, ra, dec float64, date time.Time, path string) error {
	if c.webSocketClient == nil {
		err := c.init(lon, lat, date, path)
		if err != nil {
			return err
		}
	}

	// Perform the "Start goto" request
	err := c.webSocketClient.StartGoto(ra, dec, lon, lat, date, path)
	if err != nil {
		fmt.Println("Error performing 'Start goto' request:", err)
		return err
	}

	fmt.Println("Goto operation succeeded!")

	return nil
}

func (c *dwarf2Client) UpdateDateTime() error {
	if c.webSocketClient != nil {
		err := c.webSocketClient.UpdateDateTime(time.Now())
		if err != nil {
			fmt.Println("Error updating date and time:", err)
			return err
		}
		return nil
	}

	err := errors.New("no active connection")
	fmt.Sprintf("Error updating date and time: %s", err)
	return err
}

func (c *dwarf2Client) Disconnect() error {
	if c.webSocketClient != nil {
		err := c.webSocketClient.Disconnect()
		if err != nil {
			fmt.Println("Error disconnecting from WebSocket:", err)
			return err
		}
		return nil
	}

	err := errors.New("no active connection")
	fmt.Sprintf("Error updating date and time: %s", err)
	return err
}
