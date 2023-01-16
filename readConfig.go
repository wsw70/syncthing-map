package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// readConfigXML takes the hostname and config file and dumps a data-cli.json representation for further use
func readConfigXml(hostname string, configFile string, dataFilename string) {
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
				// try to see if the names match without case
				if strings.EqualFold(hostname, knownDevice.Name) {
					log.Warn().Msgf("provided device name %s matches %s caseless", hostname, knownDevice.Name)
				}
			}
		}
	}
	if localDeviceId == "" {
		log.Fatal().Msgf("could not match the provided device name %s with known devices in the config file. Check for a warning above if this is not just a matter of case", hostname)
	}

	writeConfigToFile(fmt.Sprintf("%s %s", hostname, localDeviceId), config, dataFilename)
}

func writeConfigToFile(deviceKey string, config Configuration, dataFilename string) {
	var err error

	dataToWrite := make(dataJsonT)
	data, err := readFile(dataFilename)
	if err != nil {
		// no data file, create new one
		log.Info().Msgf("no %s file", dataFilename)
	} else {
		// unmarshall content and update with new config
		err = json.Unmarshal(data, &dataToWrite)
		if err != nil {
			log.Fatal().Msgf("cannot unmarshal %s: %s", dataFilename, err)
		}
	}
	dataToWrite[deviceKey] = config.Folder
	dataToWriteJson, err := json.Marshal(dataToWrite)
	if err != nil {
		log.Fatal().Msgf("cannot marshal data for %s: %s", dataFilename, err)
	}
	err = os.WriteFile(dataFilename, dataToWriteJson, 0644)
	if err != nil {
		log.Fatal().Msgf("cannot write %s: %s", dataFilename, err)
	}
	log.Info().Msgf("wrote %s file", dataFilename)
}
