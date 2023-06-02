package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TelescopeController struct {
	dwarfClient                *dwarf2Client
	telescopeData              *TelescopeData
	coordinatesReceivedHandler func(float64, float64)
	targetDeclination          float64
	targetRightAscension       float64
	latitude                   float64
	longitude                  float64
	utcOffsetHours             float64
	localTime                  time.Time
}

func NewTelescopeController() *TelescopeController {
	return &TelescopeController{
		dwarfClient:   &dwarf2Client{},
		telescopeData: &TelescopeData{},
	}
}

func (c *TelescopeController) GetCurrentRA() float64 {
	return c.telescopeData.RightAscension
}

func (c *TelescopeController) GetCurrentDEC() float64 {
	return c.telescopeData.Declination
}

func (c *TelescopeController) HandleCommand(command string) string {
	switch {
	case strings.HasPrefix(command, ":St"):
		degrees, err1 := strconv.ParseFloat(command[3:5], 64)
		minutes, err2 := strconv.ParseFloat(command[6:8], 64)
		if err1 == nil && err2 == nil {
			sign := 1.0
			if command[2:3] == "-" {
				sign = -1.0
			}
			c.latitude = sign * (degrees + (minutes / 60))
			return "1"
		} else {
			return "0"
		}

	case strings.HasPrefix(command, ":Sg"):
		degrees, err1 := strconv.ParseFloat(command[3:6], 64)
		minutes, err2 := strconv.ParseFloat(command[7:9], 64)
		if err1 == nil && err2 == nil {
			c.longitude = degrees + (minutes / 60)
			return "1"
		} else {
			return "0"
		}

	case strings.HasPrefix(command, ":SG"):
		utcOffsetHours, err := strconv.ParseFloat(command[3:], 64)
		if err == nil {
			c.utcOffsetHours = utcOffsetHours
			return "1"
		} else {
			return "0"
		}

	case strings.HasPrefix(command, ":SL"):
		timeStr := command[3:]
		time, err := time.Parse("15:04:05", timeStr)
		if err == nil {
			c.localTime = time
			return "1"
		} else {
			return "0"
		}

	case strings.HasPrefix(command, ":SC"):
		dateStr := command[3:]
		date, err := time.Parse("01/02/06", dateStr)
		if err == nil {
			c.localTime = time.Date(date.Year(), date.Month(), date.Day(), c.localTime.Hour(), c.localTime.Minute(), c.localTime.Second(), c.localTime.Nanosecond(), date.Location())
			return "1"
		} else {
			return "0"
		}

	case strings.HasPrefix(command, ":Sr"):
		hours, err1 := strconv.ParseFloat(command[3:5], 64)
		minutes, err2 := strconv.ParseFloat(command[6:8], 64)
		seconds, err3 := strconv.ParseFloat(command[9:11], 64)
		if err1 == nil && err2 == nil && err3 == nil {
			c.targetRightAscension = hours + ((minutes - 2) / 60) + (seconds / 3600)
			return "1"
		} else {
			return "0"
		}

	case strings.HasPrefix(command, ":Sd"):
		degrees, err1 := strconv.ParseFloat(command[4:6], 64)
		minutes, err2 := strconv.ParseFloat(command[7:9], 64)
		seconds, err3 := strconv.ParseFloat(command[10:12], 64)
		if err1 == nil && err2 == nil && err3 == nil {
			sign := 1.0
			if command[3:4] == "-" {
				sign = -1.0
			}
			c.targetDeclination = sign * (degrees + (minutes / 60) + (seconds / 3600))
			return "1"
		} else {
			return "0"

		}

	case command == ":GR":
		// Get current right ascension
		ra := c.GetCurrentRA()
		raTimeSpan := time.Duration(ra * float64(time.Hour))
		return fmt.Sprintf("+%02d:%02d.%01d#", raTimeSpan.Hours(), raTimeSpan.Minutes(), raTimeSpan.Seconds()/10)

	case command == ":MS":
		// Move telescope to specified coordinates
		if c.dwarfClient == nil {
			c.dwarfClient = &dwarf2Client{}
		}

		err := c.dwarfClient.Goto(c.longitude, c.latitude, c.targetRightAscension, c.targetDeclination, time.Now(), fmt.Sprintf("DWARF_GOTO_LX200%s", time.Now()))
		if err != nil {
			return "0"
		}

		c.telescopeData.Declination = c.targetDeclination
		c.telescopeData.RightAscension = c.targetRightAscension

		return "1"

	case command == ":GD":
		// Get current declination
		dec := c.GetCurrentDEC()
		sign := "+"
		if dec < 0 {
			sign = "-"
			dec = -dec
		}
		degrees := int(dec)
		minutes := int((dec - float64(degrees)) * 60)
		return fmt.Sprintf("%s%02d*%02d#", sign, degrees, minutes)

	case command == ":CM":
		c.telescopeData.Declination = c.targetDeclination
		c.telescopeData.RightAscension = c.targetRightAscension
		return "Coordinates matched.#"

	case strings.HasPrefix(command, ":RM"), strings.HasPrefix(command, ":RS"), strings.HasPrefix(command, ":RC"), strings.HasPrefix(command, ":RG"):
		// Initialize telescope client
		c.dwarfClient.init(c.longitude, c.latitude, time.Now(), fmt.Sprintf("DWARF_GOTO_LX200%s", time.Now()))
		return "1"

	case command == ":GVP":
		return "LX200 Custom#"

	case command == ":GVN":
		return "1.0.0#"

	case command == ":GVD":
		return "2023-05-05#"

	case command == ":GVZ":
		return "12:34:56#"
	}

	return "1"
}

func (c *TelescopeController) SetCoordinatesReceivedHandler(handler func(ra, dec float64)) {
	c.coordinatesReceivedHandler = handler
}

// UpdateTelescopeData updates the telescope data with the given right ascension and declination
func (c *TelescopeController) UpdateTelescopeData(rightAscension, declination float64) {
	c.telescopeData.RightAscension = rightAscension
	c.telescopeData.Declination = declination

	if c.coordinatesReceivedHandler != nil {
		c.coordinatesReceivedHandler(rightAscension, declination)
	}
}
