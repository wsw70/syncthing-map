package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/yaml.v2"
)

const dataFilename = "data-server.json"
const outputFilename = "syncthing-map-server.html"

// configuration file for each device and the folder holding config.xml
type serverConfigT []struct {
	Device string
	Folder string
}

func httpServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeHtmlFile()
		response, err := readFile(outputFilename)
		if err != nil {
			log.Error().Msgf("cannot read generated file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("cannot read generated file: %v", err)))
		} else {
			w.Write(response)
		}
	})
	// write output file at start
	writeHtmlFile()
	log.Info().Msg("starting server on port 3000")
	http.ListenAndServe(":3000", r)
}

func readServerConf() (conf serverConfigT) {
	configContent, err := readFile("syncthing-map-server.yaml")
	if err != nil {
		log.Fatal().Msgf("cannot read server configuration file: %v", err)
	}
	err = yaml.Unmarshal(configContent, &conf)
	if err != nil {
		log.Fatal().Msgf("cannot unmarshal server configuration file: %v", err)
	}
	return conf
}

func writeHtmlFile() {
	conf := readServerConf()
	// recreate a data file with each server
	os.Remove(dataFilename)
	for _, entry := range conf {
		readConfigXml(entry.Device, filepath.Join(entry.Folder, "config.xml"), dataFilename)
		log.Info().Msgf("parsed configuration for %s", entry.Device)
	}
	// write HTML file to serve
	writeGraph(dataFilename, outputFilename)
}
