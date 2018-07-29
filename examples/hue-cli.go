//
// hue-cli
// A program written in the Go Programming Language for the Philips Hue API.
// Copyright (C) 2018 Niels de Vos
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
//

package main

import (
	"fmt"
	"os"

	hue "github.com/collinux/GoHue"
)

func main() {
	bridges, err := hue.FindBridges()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d bridges\n", len(bridges))
	for _, bridge := range bridges {
		err = bridge.GetInfo()
		if err != nil {
			fmt.Printf("ERROR: failed to get info for bridge at %s (%s)\n", bridge.IPAddress, err)
			// fall-through, just print few details
		}
		fmt.Printf("%s\n", bridgeToString(bridge))
	}
}

func bridgeToString(bridge hue.Bridge) string {
	s := fmt.Sprintf("Bridge:\n"+
		"\tIP-address: %s",
		bridge.IPAddress)

	if bridge.Info.Device.DeviceType != "" {
		s += fmt.Sprintf("\tDevice Information:\n"+
			"\t\tDeviceType: %s\n"+
			"\t\tFriendlyName: %s\n"+
			"\t\tManufacturer: %s\n"+
			"\t\tManufacturerURL: %s\n"+
			"\t\tModelDescription: %s\n"+
			"\t\tModelName: %s\n"+
			"\t\tModelNumber: %s\n"+
			"\t\tModelURL: %s\n"+
			"\t\tSerialNumber: %s\n"+
			"\t\tUDN: %s",
			s,
			bridge.Info.Device.DeviceType,
			bridge.Info.Device.FriendlyName,
			bridge.Info.Device.Manufacturer,
			bridge.Info.Device.ManufacturerURL,
			bridge.Info.Device.ModelDescription,
			bridge.Info.Device.ModelName,
			bridge.Info.Device.ModelNumber,
			bridge.Info.Device.ModelURL,
			bridge.Info.Device.SerialNumber,
			bridge.Info.Device.UDN)
	}

	return s
}
