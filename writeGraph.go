package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// embed HTML template
var (
	//go:embed syncthing-map-template.html
	templateContents []byte
)

// TODO setup logging
func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	zerolog.New(output).With().Timestamp().Logger()
}

func writeGraph() {
	var err error
	var data dataJsonT
	var mermaidCode []string

	dataBytes, err := readFile("data.json")
	if err != nil {
		log.Fatal().Msgf("cannot read data.json file: %s", err)
	}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		log.Fatal().Msgf("cannot unmarshal data.json file: %s", err)
	}

	// build mermaid.js code
	// FIXME because of a problem with mermaid.js, all destinations are prefixed with X or x to start with a letter
	mermaidCode = append(mermaidCode, "flowchart LR")
	// get every registered device
	for deviceKey, folders := range data {
		localDeviceName := strings.Split(deviceKey, " ")[0]
		localDeviceId := strings.Split(deviceKey, " ")[1]
		// FIXME append folder name with device so that they are linked in the graph
		for i := range folders {
			folders[i].ID = fmt.Sprintf("X%s+%s", folders[i].ID, localDeviceId)
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
	f, err := os.Create("syncthing-map.html")
	if err != nil {
		log.Fatal().Msgf("cannot create syncthing-map.html: %s", err)
	}
	err = template.Execute(f, strings.Join(mermaidCode, "\n"))
	if err != nil {
		log.Fatal().Msgf("cannot populate syncthing-map.html with template: %s", err)
	}
	f.Close()
	log.Info().Msg("wrote syncthing-map.html")
}
