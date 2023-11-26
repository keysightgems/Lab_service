package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	// "strconv"
	// "net/http"
)

type Interface struct {
	Name   string `json:"name"`
	Speed  int    `json:"speed"`
	Status string `json:"status"`
	// Add more fields as needed
}
type Devicecount struct {
	Duts []Device `json:"duts"`
	Aets []Device `json:"aets"`
}

// Device represents the structure of a device
type Device struct {
	ID         float64     `json:"id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	State      string      `json:"state"`
	Attributes Attributes  `json:"attributes"`
	Interfaces []interface{} `json:"interfaces"`
	// Add more fields as needed
}

// Inventory represents the inventory details of a device
type Attributes struct {
	// Add more fields as needed
	Vendor string `json:"vendor"`
}

// Dut represents the structure of the device under "duts" key
type Dut struct {
	Name    string `json:"desc"`
	Devices Devicecount
}

// Counter represents a simple counter for generating auto-increment IDs
type Counter struct {
	Value int
}

func (c *Counter) nextID() int {
	c.Value++
	return c.Value
}

// AddDevice adds a new device with interfaces and an auto-incrementing ID to the provided map
func AddDevice(counter *Counter, devices map[int]Device, id float64, name, deviceType, state string, manufacturer string, interfaces []interface{}) map[int]Device {
	deviceID := counter.nextID()
	devices[deviceID] = Device{
		ID:    id,
		Name:  name,
		Type:  deviceType,
		State: state,
		Attributes: Attributes{
			// Add inventory details here
			Vendor: manufacturer,
		},
		Interfaces: interfaces,
	}
	return devices
}

func main() {

	cmd := exec.Command("python", "get_inventory_details.py")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing Python script:", err)
		return
	}
	// Parse JSON output into Go data structure
	var listOfDicts []map[string]interface{}
	err = json.Unmarshal(output, &listOfDicts)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	// Initialize an empty map for devices
	devices := make(map[int]Device)
	var devicesSlice []Device
	var atesSlice []Device
	for _, dict := range listOfDicts {
		// Get values using keys
		id, idExists := dict["Id"].(float64)
		name, nameExists := dict["Name"].(string)
		deviceType, deviceTypeExists := dict["DeviceType"].(string)
		manufacturer, manufacturerExists := dict["Manufacturer"].(string)
		state, stateExists := dict["State"].(string)
		interfaces, interfacesExists := dict["interfaces"].([]interface{})
		// Check if the keys exist in the dictionary
		if idExists && nameExists && deviceTypeExists && stateExists && manufacturerExists && interfacesExists {
			// Print or use the values
			fmt.Printf("ID: %f, Name: %s, Type: %s, Vendor: %s, Interface: %s, State: %s\n", id, name, deviceType, manufacturer, interfaces, state)
		} else {
			fmt.Println("Variable type not matching in the dictionary.")
		}
		
		idCounter := &Counter{}
		devices = AddDevice(idCounter, devices, id, name, deviceType, state, manufacturer, interfaces)
		// Convert the map to a slice		
		if deviceType == "DUT" {
			for _, device := range devices {
				devicesSlice = append(devicesSlice, device)
			}
		}

		// Convert the map to a slice
		if deviceType == "ATE" {
			for _, device := range devices {
				atesSlice = append(atesSlice, device)
			}
		}

	}
	
	// Create Dut with devices
	duts := Dut{
		Name:    "Inventory",
		Devices: Devicecount{Duts: devicesSlice, Aets: atesSlice},
	}

	// Marshal the Dut into JSON
	dutsJSON, err := json.MarshalIndent(duts, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Write JSON to a file
	err = ioutil.WriteFile("inventory.json", dutsJSON, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("JSON written to inventory.json")
}
