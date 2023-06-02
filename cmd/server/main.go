package main

import (
	"fmt"

	"com.github/philspins/MightyDwarf/internal"
)

func main() {
	config := internal.Config{
		IpAddress:   "127.0.0.1",
		Port:        9999,
		ComPortName: "/dev/tty.usbserial",
	}

	lx200Server := internal.NewLX200Server(config.Port)
	telescopeController := internal.NewTelescopeController()
	// telescopeController.CoordinatesReceived += async (sender, coordinates) =>
	// {
	// 	fmt.Print("RA: {coordinates.RA}, DEC: {coordinates.DEC}");
	// 	await restClient.SendRaDec(coordinates.RA, coordinates.DEC);
	// }

	var serialPortListener *internal.SerialPortListener

	serialPortListener, err := internal.NewSerialPortListener(config.ComPortName, telescopeController)
	if err != nil {
		serialPortListener.Start()
	}

	err = lx200Server.Start()
	if err != nil {
		fmt.Println("Error starting LX200 server:", err)
		return
	}
}
