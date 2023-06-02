package main

import (
	"fmt"
	"os"

	"com.github/philspins/MightyDwarf/internal"
)

func main() {
	comPortName := ""
	if len(os.Args) > 2 {
		comPortName = os.Args[2]
	}

	lx200Server := internal.NewLX200Server(9999)
	// lx200Server.TelescopeController.SetCoordinatesReceivedHandler(func(sender interface{}, coordinates internal.TelescopeData) {
	// 	fmt.Printf("RA: %f, DEC: %f\n", coordinates.RightAscension, coordinates.Declination)
	// 	//await restClient.SendRaDec(coordinates.RA, coordinates.DEC);
	// })

	var serialPortListener *internal.SerialPortListener
	if comPortName != "" {
		var err error
		serialPortListener, err = internal.NewSerialPortListener(comPortName, lx200Server.TelescopeController)
		if err != nil {
			fmt.Println("Error creating serial port listener:", err)
			return
		}
		serialPortListener.Start()
	}

	err := lx200Server.Start()
	if err != nil {
		fmt.Println("Error starting LX200 server:", err)
		return
	}
}
