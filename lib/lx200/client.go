package lx200

import (
	"fmt"
	"philspins/MightyDwarf/lib/websocket"
)

type Client struct {
}

func (c *Client) init(lon float64, lat float64, date string, path string) (int, error) {
	webSocketClient := new(websocket.WebSocketClient)
	err := webSocketClient.Connect()
	if err != nil {
		return 0, err
	}
	err = webSocketClient.UpdateDateTime(date)
	if err != nil {
		return 0, err
	}

	// Perform the "correction" request
	correctionCode, err := webSocketClient.SendCorrection(lon, lat, date, path)
	if err != nil {
		return 0, err
	}
	fmt.Println(correctionCode)
	err = webSocketClient.Disconnect()
	if err != nil {
		return 0, err
	}
	return correctionCode, nil
}
