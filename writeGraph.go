package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// embed HTML template
var (
	//go:embed syncthing-map-template.html
	templateContents []byte
)

func writeGraph(dataFilename string, outputFilename string) {
	var err error
	var data dataJsonT
	var mermaidCode []string

	dataBytes, err := readFile(dataFilename)
	if err != nil {
		log.Fatal().Msgf("cannot read %s file: %s", dataFilename, err)
	}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		log.Fatal().Msgf("cannot unmarshal %s file: %s", dataFilename, err)
	}

	// build mermaid.js code
	mermaidCode = append(mermaidCode, "flowchart LR")
	// get every registered device
	for deviceKey, folders := range data {
		localDeviceName := strings.Split(deviceKey, " ")[0]
		localDeviceId := strings.Split(deviceKey, " ")[1]
		for i := range folders {
			folders[i].ID = hashData([]string{folders[i].ID, localDeviceId})
		}
		// container for the device
		mermaidCode = append(mermaidCode, fmt.Sprintf("subgraph x%s[\"%s\"]", localDeviceId, localDeviceName))
		for _, folder := range folders {
			mermaidCode = append(mermaidCode, fmt.Sprintf("%s[\"%s\"]", folder.ID, folder.Label))
		}
		mermaidCode = append(mermaidCode, "end")
		// add shares
		for _, folder := range folders {
			for _, device := range folder.Device {
				// discard path to local device
				if device.ID == localDeviceId {
					continue
				}
				mermaidCode = append(mermaidCode, fmt.Sprintf("%s--\"%s %s\"-->x%s", folder.ID, folder.Type, device.Name, device.ID))
			}
		}
	}

	// export graph
	template, err := template.New("graph").Parse(string(templateContents))
	if err != nil {
		log.Fatal().Msgf("parse template syncthing-map-template.html: %s", err)
	}
	// var output io.Writer
	f, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal().Msgf("cannot create %s: %s", outputFilename, err)
	}
	err = template.Execute(f, strings.Join(mermaidCode, "\n"))
	if err != nil {
		log.Fatal().Msgf("cannot populate %s with template: %s", outputFilename, err)
	}
	f.Close()
	log.Info().Msgf("wrote %s", outputFilename)
}

func hashData(data []string) string {
	hsha2 := sha256.Sum256([]byte(strings.Join(data, " ")))
	return fmt.Sprintf("X%x", hsha2)
}
