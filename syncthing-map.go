package main

import (
	"io"
	"os"

	"github.com/integrii/flaggy"
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

func main() {
	// parse command line for actions
	flaggy.SetName("syncthing-map")
	flaggy.SetDescription("map syncthing devices and folders")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	// read config file & hostname
	readSub := flaggy.NewSubcommand("read")
	var hostname string
	var configFile string
	readSub.Description = "read config file"
	readSub.String(&hostname, "d", "device", "hostname for the config file")
	readSub.String(&configFile, "f", "file", "XML config file")
	flaggy.AttachSubcommand(readSub, 1)
	// create graph
	graphSub := flaggy.NewSubcommand("graph")
	flaggy.AttachSubcommand(graphSub, 2)
	// clean up
	cleanSub := flaggy.NewSubcommand("clean")
	flaggy.AttachSubcommand(cleanSub, 3)
	// parse the command line
	flaggy.Parse()
	if readSub.Used {
		readConfigXml(hostname, configFile)
	} else if graphSub.Used {
		writeGraph()
	} else if cleanSub.Used {
		//
	} else {
		flaggy.ShowHelpAndExit("🛑 missing command")
	}
}
