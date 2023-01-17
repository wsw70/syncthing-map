package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

// version will be added from tag during compilation
var compiledVersion string

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
	if compiledVersion == "" {
		log.Error().Msgf("compiledVersion not set at compile time")
		compiledVersion = "(missing from compilation)"
	}

	app := &cli.App{
		Name:    "syncthing-map",
		Usage:   "Syncthing devices and shared folders mapped in your browser",
		Version: fmt.Sprintf("%s %s/%s", compiledVersion, runtime.GOOS, runtime.GOARCH),
		Authors: []*cli.Author{
			{
				Name:  "wsw70",
				Email: "1345886+wsw70@users.noreply.github.com",
			},
		},
		Copyright: "WTFPL (http://www.wtfpl.net)",
		HideHelp:  true,
		UsageText: "syncthing-map clean\nsyncthing-map add --device <device name> --file <configuration file> | graph\nsyncthing-map server",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a new config file",
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("device") == "" || cCtx.String("file") == "" {
						cli.ShowAppHelpAndExit(cCtx, 1)
					}
					readConfigXml(cCtx.String("device"), cCtx.String("file"), "data-cli.json")
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
				Usage:   "create the graph in syncthing-map.html",
				Action: func(cCtx *cli.Context) error {
					writeGraph("data-cli.json", "syncthing-map-cli.html")
					return nil
				},
			},
			{
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "remove working files (*.json, *.html)",
				Action: func(cCtx *cli.Context) error {
					filesToRemove := []string{
						"syncthing-map-cli.html",
						"syncthing-map-server.html",
						"data-cli.json",
						"data-server.json",
					}
					for _, file := range filesToRemove {
						os.Remove(file)
					}
					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "HTTP server to automatically generate a map",
				Action: func(cCtx *cli.Context) error {
					httpServer()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msgf("error running the application: %v", err)
	}
}
