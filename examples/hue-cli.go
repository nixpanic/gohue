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
	"errors"
	"fmt"
	"os"

	hue "github.com/collinux/GoHue"
	"github.com/spf13/cobra"

	// for yaml conversion of the ConfigFile
	"gopkg.in/yaml.v2"
)

var (
	// `hue-cli` parameters
	bridgeIP   string
	configFile string

	// `hue-cli create-user` parameters
	deviceName string
)

type ConfigFile struct {
	Bridges []BridgeConfig
}

type BridgeConfig struct {
	IPAddress string `yaml:"ipaddress"`
	User      string `yaml:"user"`
}

func init() {
	// hue-cli --config=<filename>
	cmdHueCli.Flags().StringVar(&configFile, "config", "",
		"filename with the configuration (optional)")

	// hue-cli discover
	cmdHueCli.AddCommand(cmdDiscover)
	// hue-cli discover --bridge=<ip-address>
	cmdDiscover.Flags().StringVar(&bridgeIP, "bridge", "",
		"IP-address of the bridge (optional)")

	// hue-cli create-user
	cmdHueCli.AddCommand(cmdCreateUser)
	// hue-cli create-user --bridge=<ip-address>
	cmdCreateUser.Flags().StringVar(&bridgeIP, "bridge", "",
		"IP-address of the bridge (optional)")
	// hue-cli create-user --device=<name>
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	cmdCreateUser.Flags().StringVar(&deviceName, "device", hostname,
		"name of the device hue-cli is running on (optional)")
}

var cmdHueCli = &cobra.Command{
	Use:   "hue-cli",
	Short: "Commandline application to show the capabilities of GoHue",
	Long:  "Commandline application to show the capabilities of GoHue",
}

var cmdDiscover = &cobra.Command{
	Use:   "discover",
	Short: "discover bridges",
	Long:  "request the known bridges in this network from https://discovery.meethue.com/",

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var bridges []hue.Bridge

		if bridgeIP != "" {
			// if we know the IP-addres, we dont do any discovery
			bridge, err := hue.NewBridge(bridgeIP)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to find bridge %s: %s", bridgeIP, err))
			}

			bridges = append(bridges, *bridge)
		} else {
			bridges, err = hue.FindBridges()
			if err != nil {
				return err
			}
		}

		fmt.Printf("Found %d bridges\n", len(bridges))
		for _, bridge := range bridges {
			err := bridge.GetInfo()
			if err != nil {
				fmt.Printf("ERROR: failed to get info for bridge at %s (%s)\n", bridge.IPAddress, err)
				// fall-through, just print few details
			}
			fmt.Printf("%s\n", bridgeToString(bridge))
		}
		return nil
	},
}

var cmdCreateUser = &cobra.Command{
	Use:   "create-user",
	Short: "create a new user on the bridge",
	Long:  "create a new user on the bridge, should have pressed the 'link button' in advance",

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var bridge *hue.Bridge

		// start of bridge detection (TODO: move to its own function?)
		if bridgeIP != "" {
			bridge, err = hue.NewBridge(bridgeIP)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to find bridge %s: %s", bridgeIP, err))
			}
		} else {
			// if no bridge is given, detect them all
			// in case more than one bridge is found, error out
			bridges, err := hue.FindBridges()
			if err != nil {
				return err
			} else if len(bridges) == 0 {
				return errors.New("no bridge found")
			} else if len(bridges) > 1 {
				return errors.New(fmt.Sprintf("%d bridges found, use --bridge=<ip-address>", len(bridges)))
			}

			bridge = &bridges[0]
		}

		// we got a bridge, create a new user
		user, err := bridge.CreateUser("hue-cli#"+deviceName)
		if err != nil {
			return err
		}

		// generate a new config file
		config := &ConfigFile{
			Bridges: []BridgeConfig{{
				IPAddress: bridge.IPAddress,
				User: user,
			}},
		}

		configYaml, err := yaml.Marshal(&config)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to conver to yaml (%s)", err))
		}

		fmt.Printf("new configuration: %s\n", configYaml)

		return nil
	},
}

func main() {
	err := cmdHueCli.Execute()
	if err != nil {
		os.Exit(1)
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
