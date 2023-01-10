package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

// readConfigXML takes the hostname and config file and dumps a data.json representation for further use
func readConfigXml(hostname string, configFile string) {
	var err error
	var localDeviceId string

	config := Configuration{}
	fileContent, err := readFile(configFile)
	if err != nil {
		log.Fatal().Msgf("cannot read %s: %s", configFile, err)
	}
	err = xml.Unmarshal(fileContent, &config)
	if err != nil {
		log.Fatal().Msgf("cannot unmarshall config.xml: %v", err)
	}
	// update the devices names in folders
	for _, folder := range config.Folder {
		for i, device := range folder.Device {
			// find the appropraite device name
			for _, knownDevice := range config.Device {
				// found matching device ID in devices
				if device.ID == knownDevice.ID {
					folder.Device[i].Name = knownDevice.Name
				}
				// attempt to set the ID of the device if names match
				if hostname == knownDevice.Name {
					localDeviceId = knownDevice.ID
				}
			}
		}
	}

	writeConfigToFile(fmt.Sprintf("%s %s", hostname, localDeviceId), config)
}

// format of data of data.json
const dataFile = "data.json"

func writeConfigToFile(deviceKey string, config Configuration) {
	var err error

	dataToWrite := make(dataJsonT)
	data, err := readFile(dataFile)
	if err != nil {
		// no data file, create new one
		log.Info().Msgf("no %s file", dataFile)
	} else {
		// unmarshall content and update with new config
		err = json.Unmarshal(data, &dataToWrite)
		if err != nil {
			log.Fatal().Msgf("cannot unmarshal %s: %s", dataFile, err)
		}
	}
	dataToWrite[deviceKey] = config.Folder
	dataToWriteJson, err := json.Marshal(dataToWrite)
	if err != nil {
		log.Fatal().Msgf("cannot marshal data for %s: %s", dataFile, err)
	}
	err = os.WriteFile(dataFile, dataToWriteJson, 0644)
	if err != nil {
		log.Fatal().Msgf("cannot write %s: %s", dataFile, err)
	}
	log.Info().Msgf("wrote data.json file")
}
