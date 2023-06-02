package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type websock struct {
	Config          Config
	clientWebSocket *websocket.Conn
}

func (c *websock) GetServerUrl(protocol string, port int) string {
	return fmt.Sprintf("%s://%s:%d", protocol, c.Config.IpAddress, port)
}

func (c *websock) Connect() error {
	serverURI := c.GetServerUrl("ws", 9900)
	wsHeaders := http.Header{}
	wsHeaders.Add("Origin", serverURI)

	conn, _, err := websocket.DefaultDialer.Dial(serverURI, wsHeaders)
	if err != nil {
		return err
	}

	c.clientWebSocket = conn
	return nil
}

func (c *websock) SendMessage(message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = c.clientWebSocket.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (c *websock) ReceiveMessage(message interface{}) error {
	_, jsonData, err := c.clientWebSocket.ReadMessage()
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, message)
	if err != nil {
		return err
	}

	return nil
}

func (c *websock) Disconnect() error {
	err := c.clientWebSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	err = c.clientWebSocket.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *websock) TurnOnCamera() error {
	request := map[string]interface{}{
		"Interface": 10000,
		"CamId":     0,
	}

	err := c.SendMessage(request)
	if err != nil {
		return err
	}

	return nil
}

func (c *websock) UpdateDateTime(date time.Time) error {
	dateTime := date.Format("2006-01-02 15:04:05")
	path := fmt.Sprintf("%s/date?date=%s", c.GetServerUrl("http", 8092), url.QueryEscape(dateTime))

	response, err := http.Get(path)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update date and time: %s", response.Status)
	}

	return nil
}

func (c *websock) SendCorrection(lon, lat float64, date time.Time, path string) error {
	request := map[string]interface{}{
		"Interface": 11205,
		"CamId":     0,
		"Lon":       lon,
		"Lat":       lat,
		"Date":      date,
		"Path":      path,
	}

	err := c.SendMessage(request)
	if err != nil {
		return err
	}

	var message []byte
	err = c.ReceiveMessage(&message)
	if err != nil {
		return err
	}

	return nil
}

func (c *websock) StartGoto(ra, dec, lon, lat float64, date time.Time, path string) error {
	request := map[string]interface{}{
		"Interface": 11203,
		"CamId":     0,
		"Ra":        ra,
		"Dec":       dec,
		"Lon":       lon,
		"Lat":       lat,
		"Date":      date,
		"Path":      path,
	}

	err := c.SendMessage(request)
	if err != nil {
		return err
	}

	var message []byte
	err = c.ReceiveMessage(&message)
	if err != nil {
		return err
	}

	return nil
}
