package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

type Device struct {
	ID   string `xml:"id,attr" json:"id"`
	Name string `xml:"name,attr" json:"name"`
}

type Folder struct {
	ID     string   `xml:"id,attr" json:"id"`
	Label  string   `xml:"label,attr" json:"label"`
	Type   string   `xml:"type,attr" json:"type"`
	Device []Device `xml:"device" json:"device"`
}

type Configuration struct {
	Folder []Folder `xml:"folder"`
	Device []Device `xml:"device"`
}

type dataJsonT map[string][]Folder

func readFile(filename string) (content []byte, err error) {
	handler, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	content, err = io.ReadAll(handler)
	if err != nil {
		return nil, err
	}
	handler.Close()
	return content, nil
}

// global setting for data files
const dataFilename = "data.json"
const outputFilename = "syncthing-map.html"

var log zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	log = zerolog.New(output).With().Timestamp().Logger()
}

func main() {

	app := &cli.App{
		Name:     "syncthing-map",
		Version:  "alpha",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "wsw70",
				Email: "1345886+wsw70@users.noreply.github.com",
			},
		},
		Copyright: "WTFPL http://www.wtfpl.net/",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a new config file",
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("device") == "" || cCtx.String("file") == "" {
						cli.ShowAppHelpAndExit(cCtx, 1)
					}
					readConfigXml(cCtx.String("device"), cCtx.String("file"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "device", Aliases: []string{"d"}},
					&cli.StringFlag{Name: "file", Aliases: []string{"f"}},
				},
			},
			{
				Name:    "graph",
				Aliases: []string{"g"},
				Usage:   "create th egraph in syncthing-map.html",
				Action: func(cCtx *cli.Context) error {
					writeGraph()
					return nil
				},
			},
			{
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "remove working files (data.json, syncthing-map.html)",
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msgf("error running the application: %v", err)
	}
}
